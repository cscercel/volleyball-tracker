import streamlit as st
from datetime import datetime
from api import api
from utils.court import plot_match_court
from utils.misc_functions import shuffle_players


st.set_page_config(page_title="Matches", page_icon="🏐", layout="wide")

st.title("🏐 Matches")

current_year = datetime.utcnow().year

tab1, tab2, tab3 = st.tabs(["Create Match", "Draft Matches", "Completed Matches"])


# ---------------------------------------------------------------------------
# TAB 1 — Create Match
# ---------------------------------------------------------------------------
with tab1:
    st.subheader("Create New Match")

    try:
        players     = api.get_players()
        player_names = [p["name"] for p in players]
    except Exception:
        st.error("Failed to load players.")
        player_names = []

    if len(player_names) < 2:
        st.warning("⚠️ Need at least 2 players to create a match. Add players first!")
    else:
        col1, col2 = st.columns(2)

        with col1:
            match_type = st.pills("Match Type", ["indoor", "beach"], default="indoor", key="create_type")
            blue_team  = st.multiselect("Select Blue Team", player_names, key="blue")

        with col2:
            draft_type = st.pills("Draft Type", ["random", "manual"], default="random", key="draft_type")
            available_for_red = [p for p in player_names if p not in blue_team]
            red_team = st.multiselect("Select Red Team", available_for_red, key="red")

        if st.button("Create Match", type="primary"):
            if not blue_team or not red_team:
                st.error("Both teams need at least one player!")
            else:
                if draft_type == "random":
                    blue_team, red_team = shuffle_players(blue_team + red_team)

                try:
                    match = api.create_match({
                        "match_type": match_type,
                        "blue_team":  blue_team,
                        "red_team":   red_team,
                    })
                    st.success("✅ Match created!")

                    # Fetch roster to show team breakdown
                    try:
                        roster = api.get_match_roster(match["id"])
                        blue_players = [r["player_name"] for r in roster if r["color"] == "blue"]
                        red_players  = [r["player_name"] for r in roster if r["color"] == "red"]
                    except Exception:
                        blue_players = blue_team
                        red_players  = red_team

                    blue_col, red_col = st.columns(2, border=True)
                    with blue_col:
                        st.subheader("🔵 Blue Team")
                        st.markdown(", ".join(blue_players))
                    with red_col:
                        st.subheader("🔴 Red Team")
                        st.markdown(", ".join(red_players))

                    fig = plot_match_court(blue_players, red_players)
                    st.pyplot(fig)

                except Exception as e:
                    st.error(f"Failed to create match: {e}")


# ---------------------------------------------------------------------------
# TAB 2 — Draft (uncompleted) Matches
# ---------------------------------------------------------------------------
with tab2:
    st.subheader("Draft Matches (Not Played Yet)")

    try:
        drafts = api.get_uncompleted_matches()
    except Exception:
        st.error("Failed to load draft matches.")
        drafts = []

    if not drafts:
        st.info("No draft matches. Create one in the 'Create Match' tab!")
    else:
        for match in drafts:
            match_id = match["id"]

            # Fetch roster for each draft match
            try:
                roster = api.get_match_roster(match_id)
                blue_players = [r["player_name"] for r in roster if r["color"] == "blue"]
                red_players  = [r["player_name"] for r in roster if r["color"] == "red"]
            except Exception:
                blue_players, red_players = [], []

            with st.expander(f"Match {match_id[:8]} — {match['match_type'].upper()}"):
                col1, col2 = st.columns(2)
                with col1:
                    st.markdown("**🔵 Blue Team**")
                    for name in blue_players:
                        st.text(f"- {name}")
                with col2:
                    st.markdown("**🔴 Red Team**")
                    for name in red_players:
                        st.text(f"- {name}")

                st.markdown("---")
                st.markdown("**Submit Results**")

                col_a, col_b, col_c = st.columns(3)
                with col_a:
                    blue_score = st.number_input(
                        "Blue Score", min_value=None, value=None, key=f"blue_{match_id}"
                    )
                with col_b:
                    red_score = st.number_input(
                        "Red Score", min_value=None, value=None, key=f"red_{match_id}"
                    )
                with col_c:
                    if st.button("Submit Results", key=f"submit_{match_id}"):
                        if blue_score is None or red_score is None:
                            st.error("Enter scores for both teams.")
                        elif blue_score == red_score:
                            st.error("Scores cannot be equal — a winner must be determined.")
                        else:
                            try:
                                api.submit_match_results(match_id, int(blue_score), int(red_score))
                                st.success("✅ Results submitted!")
                                st.rerun()
                            except Exception:
                                st.error("Failed to submit results.")

                    if st.button("Delete Draft", key=f"delete_{match_id}"):
                        try:
                            api.delete_match(match_id)
                            st.success("✅ Draft deleted!")
                            st.rerun()
                        except Exception:
                            st.error("Failed to delete draft.")


# ---------------------------------------------------------------------------
# TAB 3 — Completed Matches
# ---------------------------------------------------------------------------
with tab3:
    st.subheader("Completed Matches")

    col1, col2 = st.columns(2)
    with col1:
        completed_type = st.selectbox("Match Type", ["indoor", "beach"], key="completed_type")
    with col2:
        completed_season = st.number_input(
            "Season", min_value=2023, max_value=current_year, value=current_year, key="completed_season"
        )

    try:
        completed = api.get_matches_by_season(completed_type, completed_season)
    except Exception:
        st.error("Failed to load completed matches.")
        completed = []

    if not completed:
        st.info("No completed matches for this filter.")
    else:
        for match in completed:
            match_id   = match["id"]
            blue_score = match["blue_score"]
            red_score  = match["red_score"]

            # Determine winner & OT from scores (backend doesn't return winner field)
            is_otl = abs(blue_score - red_score) == 2
            winner = "blue" if blue_score > red_score else "red"
            winner_emoji = "🔵" if winner == "blue" else "🔴"
            ot_badge = " ⏱️ OT" if is_otl else ""

            created = match.get("created_at", "")[:10]

            with st.expander(
                f"{winner_emoji} {blue_score}–{red_score}{ot_badge} — {completed_type.upper()} — {created}"
            ):
                # Fetch roster
                try:
                    roster = api.get_match_roster(match_id)
                    blue_players = [r["player_name"] for r in roster if r["color"] == "blue"]
                    red_players  = [r["player_name"] for r in roster if r["color"] == "red"]
                except Exception:
                    blue_players, red_players = [], []

                col1, col2 = st.columns(2)
                with col1:
                    st.markdown(f"**🔵 Blue Team — {blue_score}**")
                    for name in blue_players:
                        st.text(f"- {name}")
                with col2:
                    st.markdown(f"**🔴 Red Team — {red_score}**")
                    for name in red_players:
                        st.text(f"- {name}")
