const RANKS: [string, number, number][] = [
    ['Iron I', 0, 0.1],
    ['Iron II', 0.1, 0.2],
    ['Iron III', 0.2, 0.3],
    ['Bronze I', 0.3, 0.4],
    ['Bronze II', 0.4, 0.5],
    ['Bronze III', 0.5, 0.6],
    ['Silver I', 0.6, 0.7],
    ['Silver II', 0.7, 0.8],
    ['Silver III', 0.8, 0.9],
    ['Gold I', 0.9, 1.0],
    ['Gold II', 1.0, 1.1],
    ['Gold III', 1.1, 1.2],
    ['Platinum I', 1.2, 1.3],
    ['Platinum II', 1.3, 1.4],
    ['Platinum III', 1.4, 1.5],
    ['Diamond I', 1.5, 1.6],
    ['Diamond II', 1.6, 1.7],
    ['Diamond III', 1.7, 1.8],
    ['Spiker', 1.8, 1.9],
    ['Ace', 1.9, 2.0],
    ['Sensei', 2.0, Infinity],
]

export function calculateMmr(avgPoints: number, efficiencyRate: number): number {
    return avgPoints * efficiencyRate
}

export function getRank(played: number, points: number, efficiencyRate: number): string {
    if (played < 10) return 'Unranked'
    const avgPoints = played > 0 ? points / played : 0
    const mmr = calculateMmr(avgPoints, efficiencyRate)
    for (const [name, low, high] of RANKS) {
        if (mmr >= low && mmr < high) return name
    }
    return 'Iron I'
}
