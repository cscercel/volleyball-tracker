import math
import streamlit as st
from datetime import datetime
from api import api
from utils.misc_functions import calculate_mmr, get_rank


st.set_page_config(page_title="Players", page_icon="👥", layout="wide")

st.title("👥 Players")

current_year = datetime.utcnow().year

tab1, tab2, tab3 = st.tabs(["Player Profile", "Add Player", "Manage Players"])


# ---------------------------------------------------------------------------
# TAB 1 — Player Profile
# ---------------------------------------------------------------------------
with tab1:
    st.subheader("Player Stats")

    try:
        roster = api.get_players()
    except Exception:
        roster = []

    if not roster:
        st.info("No players yet!")
        st.stop()

    # Build name → id lookup
    name_to_id = {p["name"]: p["id"] for p in roster}
    name_list = list(name_to_id.keys())

    selected_name = st.selectbox("Player", name_list)
    player_id = name_to_id[selected_name]

    col_type, col_season = st.columns(2)
    with col_type:
        match_type = st.selectbox("Match Type", ["indoor", "beach"], key="profile_type")
    with col_season:
        season = st.number_input(
            "Season", min_value=2023, max_value=current_year, value=current_year, key="profile_season"
        )

    # ---- Stats ----
    try:
        stats = api.get_player_stats(player_id, match_type, season)
        has_stats = True
    except Exception:
        has_stats = False

    try:
        prev_stats = api.get_player_stats(player_id, match_type, season - 1)
        has_prev = True
    except Exception:
        has_prev = False

    st.markdown("---")

    if not has_stats:
        st.info(f"No {match_type} stats for season {season}.")
    else:
        played         = stats["played"]
        wins           = stats["wins"]
        losses         = stats["losses"]
        otl            = stats["otl"]
        points         = stats["points"]
        streak         = stats["streak"]
        longest_streak = stats["longest_streak"]
        win_rate       = float(stats.get("win_rate") or 0)
        efficiency     = float(stats.get("efficiency_rate") or 0)
        avg_points     = points / played if played > 0 else 0
        mmr            = calculate_mmr(avg_points, efficiency)
        rank           = get_rank(mmr, played)

        c1, c2, c3 = st.columns(3)
        with c1:
            delta_played = (played - has_prev * prev_stats["played"]) if has_prev else None
            st.metric("Matches Played", played, delta=delta_played)
            st.metric("Current Win Streak", streak)
        with c2:
            wr_pct = round(win_rate * 100, 2)
            delta_wr = round(wr_pct - float(prev_stats.get("win_rate") or 0) * 100, 2) if has_prev else None
            st.metric("Win Rate", f"{wr_pct}%", delta=delta_wr)
            delta_ls = (longest_streak - prev_stats["longest_streak"]) if has_prev else None
            st.metric("Longest Streak", longest_streak, delta=delta_ls)
        with c3:
            delta_pts = (points - prev_stats["points"]) if has_prev else None
            st.metric("Points", points, delta=delta_pts)
            st.metric("W / L / OTL", f"{wins} / {losses} / {otl}")

        # Rank badge
        st.image(f"assets/{rank}.png", caption=rank, width=160)
        if rank == "Unranked":
            st.markdown(f"Play **{10 - played}** more games to receive a rank!")
        elif rank == "Sensei":
            st.markdown("_Through Heaven and Earth I alone am honored_")
        else:
            mmr_progress = mmr - math.floor(mmr)
            st.progress(mmr_progress, text="Progress towards next rank")
            st.caption(f"{int(mmr_progress * 100)} / 100")

    # ---- Match History ----
    st.markdown("---")
    st.subheader("Match History")
    try:
        history = api.get_player_history(player_id, match_type, season)
        if not history:
            st.info("No completed matches this season.")
        else:
            for m in history:
                blue_score = m["blue_score"]
                red_score  = m["red_score"]
                color      = m["color"]  # "blue" or "red"

                my_score    = blue_score if color == "blue" else red_score
                their_score = red_score  if color == "blue" else blue_score
                won = my_score > their_score
                is_otl = abs(blue_score - red_score) == 2 and not won

                result_label = "✅ Win" if won else ("⚠️ OTL" if is_otl else "❌ Loss")
                created = m.get("created_at", "")
                st.write(f"{result_label} — {my_score}:{their_score} ({color} team) — {created[:10]}")
    except Exception:
        st.error("Failed to load match history.")


# ---------------------------------------------------------------------------
# TAB 2 — Add Player
# ---------------------------------------------------------------------------
with tab2:
    st.subheader("Add New Player")
    with st.form("add_player"):
        player_name = st.text_input("Player Name")
        submit = st.form_submit_button("Add Player")

        if submit and player_name:
            try:
                api.create_player(player_name)
                st.success(f"✅ Added {player_name}")
            except Exception:
                st.error("Failed to add player.")


# ---------------------------------------------------------------------------
# TAB 3 — Manage Players (rename / delete)
# ---------------------------------------------------------------------------
with tab3:
    st.subheader("Manage Players")

    if not roster:
        st.info("No players to manage.")
    else:
        manage_name = st.selectbox("Select player", name_list, key="manage_select")
        manage_id   = name_to_id[manage_name]

        with st.expander("✏️ Rename Player"):
            with st.form("rename_player"):
                new_name = st.text_input("New name", value=manage_name)
                if st.form_submit_button("Rename"):
                    try:
                        api.update_player_name(manage_id, new_name)
                        st.success(f"✅ Renamed to {new_name}")
                        st.rerun()
                    except Exception:
                        st.error("Failed to rename player.")

        with st.expander("🗑️ Delete Player"):
            st.warning(f"This will permanently delete **{manage_name}**.")
            if st.button("Confirm Delete", type="primary"):
                if not st.session_state.get("authenticated"):
                    st.error("You must be logged in to delete players.")
                else:
                    try:
                        api.delete_player(manage_id)
                        st.success(f"✅ Deleted {manage_name}")
                        st.rerun()
                    except Exception:
                        st.error("Failed to delete player.")
