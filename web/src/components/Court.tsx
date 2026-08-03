import { useMemo } from 'react'

type CourtProps = {
  blueTeam: string[]
  redTeam: string[]
}

// Court dimensions
const H = 270
const W = 450
const scaleX = W / 9
const scaleY = H / 18

// Blue positions (left side) — exact coordinates from original Python tool
const blueCoords: Record<number, [number, number]> = {
  1: [1, 2],
  2: [3.5, 2],
  3: [3.5, 9],
  4: [3.5, 16],
  5: [1, 16],
  6: [1, 9],
}

// Red positions (right side)
const redCoords: Record<number, [number, number]> = {
  7: [8, 16],
  8: [5.5, 16],
  9: [5.5, 9],
  10: [5.5, 2],
  11: [8, 2],
  12: [8, 9],
}

const positionSets: Record<number, number[]> = {
  1: [1],
  2: [1, 3],
  3: [1, 3, 5],
  4: [1, 2, 3, 5],
  5: [1, 2, 3, 4, 5],
  6: [1, 2, 3, 4, 5, 6],
}

function shuffle<T>(arr: T[]): T[] {
  return [...arr].sort(() => Math.random() - 0.5)
}

function assignTeam(team: string[], side: 'blue' | 'red') {
  const shuffled = shuffle(team)
  const onCourt = shuffled.slice(0, 6)
  const bench = shuffled.slice(6)
  const count = Math.min(onCourt.length, 6)
  const positions = positionSets[count] ?? []
  const coords = side === 'blue' ? blueCoords : redCoords
  const offset = side === 'red' ? 6 : 0

  return {
    players: onCourt.map((name, i) => {
      const pos = positions[i] + offset
      const [cx, cy] = coords[pos]
      return { name, cx: cx * scaleX, cy: cy * scaleY }
    }),
    bench,
  }
}

export default function Court({ blueTeam, redTeam }: CourtProps) {
  // Recompute only when the team rosters actually change, mirroring the
  // one-shot $derived behavior from the Svelte version.
  const { blue, red, ballX, ballY } = useMemo(() => {
    const blue = assignTeam(blueTeam, 'blue')
    const red = assignTeam(redTeam, 'red')
    const firstServe = Math.random() < 0.5 ? 'blue' : 'red'
    const ballX = firstServe === 'blue' ? 0.2 * scaleX : 8.8 * scaleX
    const ballY = firstServe === 'blue' ? 2 * scaleY : 16 * scaleY
    return { blue, red, ballX, ballY }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [blueTeam.join(','), redTeam.join(',')])

  const benchY = H + 30

  return (
    <div className="mt-6 flex justify-center">
      <svg
        viewBox={`-40 -20 ${W + 80} ${H + 70}`}
        xmlns="http://www.w3.org/2000/svg"
        className="w-full max-w-[500px] rounded-lg"
      >
        <rect x="-40" y="-20" width={W + 80} height={H + 70} fill="bisque" />
        <rect x="0" y="0" width={W} height={H} fill="bisque" stroke="darkred" strokeWidth={3} />
        <line
          x1={W / 2}
          y1={0}
          x2={W / 2}
          y2={H}
          stroke="black"
          strokeWidth={3}
          strokeDasharray="10,6"
        />
        <line
          x1={3 * scaleX}
          y1={0}
          x2={3 * scaleX}
          y2={H}
          stroke="darkred"
          strokeWidth={2}
          strokeDasharray="8,5"
          opacity={0.5}
        />
        <line
          x1={6 * scaleX}
          y1={0}
          x2={6 * scaleX}
          y2={H}
          stroke="darkred"
          strokeWidth={2}
          strokeDasharray="8,5"
          opacity={0.5}
        />

        {blue.players.map((player) => (
          <g key={`blue-${player.name}`}>
            <circle cx={player.cx} cy={player.cy} r={18} fill="blue" opacity={0.25} />
            <text
              x={player.cx}
              y={player.cy}
              textAnchor="middle"
              dominantBaseline="middle"
              fontSize={11}
              fontWeight="bold"
              fill="black"
            >
              {player.name}
            </text>
          </g>
        ))}

        {red.players.map((player) => (
          <g key={`red-${player.name}`}>
            <circle cx={player.cx} cy={player.cy} r={18} fill="red" opacity={0.25} />
            <text
              x={player.cx}
              y={player.cy}
              textAnchor="middle"
              dominantBaseline="middle"
              fontSize={11}
              fontWeight="bold"
              fill="black"
            >
              {player.name}
            </text>
          </g>
        ))}

        <circle cx={ballX} cy={ballY} r={10} fill="yellow" stroke="darkblue" strokeWidth={2} />

        {blue.bench.map((name, i) => (
          <text
            key={`blue-bench-${name}`}
            x={(i / Math.max(blue.bench.length - 1, 1)) * (W * 0.35)}
            y={benchY}
            textAnchor="middle"
            dominantBaseline="middle"
            fontSize={11}
            fontWeight="bold"
            fill="blue"
          >
            {name}
          </text>
        ))}

        {red.bench.map((name, i) => (
          <text
            key={`red-bench-${name}`}
            x={W - (i / Math.max(red.bench.length - 1, 1)) * (W * 0.35)}
            y={benchY}
            textAnchor="middle"
            dominantBaseline="middle"
            fontSize={11}
            fontWeight="bold"
            fill="red"
          >
            {name}
          </text>
        ))}

        {(blue.bench.length > 0 || red.bench.length > 0) && (
          <text
            x={W / 2}
            y={benchY}
            textAnchor="middle"
            dominantBaseline="middle"
            fontSize={10}
            fill="#666"
          >
            — bench —
          </text>
        )}
      </svg>
    </div>
  )
}
