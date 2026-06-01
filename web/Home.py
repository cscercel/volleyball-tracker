import streamlit as st
import pandas as pd
from datetime import datetime
from api import api
from utils.misc_functions import calculate_mmr, get_rank


st.set_page_config(
    page_title="Volleyball Tracker",
    page_icon="🏐",
    layout="wide",
)

st.title("🏐 Volleyball Tracker")
st.markdown("---")

st.markdown("""
### Welcome to the Volleyball Tracker!

Use the sidebar to navigate:
- **Players** - Manage your player roster
- **Matches** - Create matches and submit results
- **Leaderboard** - View stats and rankings
""")

st.markdown("---")
st.header("🏆 Leaderboard")

col1, col2 = st.columns(2)
with col1:
    match_type = st.selectbox("Match Type", ["indoor", "beach"])
with col2:
    current_year = datetime.utcnow().year
    season = st.number_input("Season", min_value=2023, max_value=current_year, value=current_year)

try:
    leaderboard = api.get_leaderboard(match_type, season)

    if not leaderboard:
        st.info(f"No stats for {match_type} in season {season}.")
    else:
        rows = []
        for entry in leaderboard:
            win_rate = float(entry.get("win_rate") or 0)
            efficiency = float(entry.get("efficiency_rate") or 0)
            avg_points = entry["points"] / entry["played"] if entry["played"] > 0 else 0
            mmr = calculate_mmr(avg_points, efficiency)
            rank = get_rank(mmr, entry["played"])
            rows.append({
                "Player":    entry["name"],
                "Played":    entry["played"],
                "Wins":      entry["wins"],
                "Losses":    entry["losses"],
                "OTL":       entry["otl"],
                "Points":    entry["points"],
                "Win Rate":  f"{win_rate:.1%}",
                "Rank":      rank,
                # Hidden helpers (excluded from display)
                "_mmr":      mmr,
                "_eff":      efficiency,
            })

        df = pd.DataFrame(rows)
        df = df.sort_values("Points", ascending=False).reset_index(drop=True)
        df.index += 1

        st.dataframe(
            df,
            use_container_width=True,
            column_config={"_mmr": None, "_eff": None},
        )

except Exception:
    st.error("Failed to load leaderboard.")
