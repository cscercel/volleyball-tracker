<script lang="ts">
// Props — equivalent to plot_match_court(blue_team, red_team)
let { blueTeam, redTeam }: { blueTeam: string[], redTeam: string[] } = $props()

// Court dimensions
const H = 270
const W = 450 
const scaleX = W / 9
const scaleY = H / 18

// Blue positions (left side) — exact coordinates from Python
const blueCoords: Record<number, [number, number]> = {
1: [1,   2],
   2: [3.5, 2],
   3: [3.5, 9],
   4: [3.5, 16],
   5: [1,   16],
   6: [1,   9],
}

// Red positions (right side)
const redCoords: Record<number, [number, number]> = {
7:  [8,   16],
    8:  [5.5, 16],
    9:  [5.5, 9],
    10: [5.5, 2],
    11: [8,   2],
    12: [8,   9],
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
        const count = Math.min(onCourt.length, 6) as 1|2|3|4|5|6
        const positions = positionSets[count]
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

const blue = $derived(assignTeam(blueTeam, 'blue'))
const red  = $derived(assignTeam(redTeam, 'red'))

// Random first serve
const firstServe = Math.random() < 0.5 ? 'blue' : 'red'
const ballX = firstServe === 'blue' ? 0.2 * scaleX : 8.8 * scaleX
const ballY = firstServe === 'blue' ? 2   * scaleY : 16  * scaleY

// Bench Y position (below court)
const benchY = H + 30
</script>

<div class="court-wrapper">
<svg
viewBox="-40 -20 {W + 80} {H + 70}"
xmlns="http://www.w3.org/2000/svg"
class="court-svg"
>
<!-- Background -->
<rect x="-40" y="-20" width={W + 80} height={H + 70} fill="bisque" />

<!-- Court boundary -->
<rect x="0" y="0" width={W} height={H} fill="bisque" stroke="darkred" stroke-width="3" />

<!-- Center line (dashed) -->
<line
x1={W / 2} y1="0"
x2={W / 2} y2={H}
stroke="black" stroke-width="3"
stroke-dasharray="10,6"
/>

<!-- Attack lines -->
<line
x1={3 * scaleX} y1="0"
x2={3 * scaleX} y2={H}
stroke="darkred" stroke-width="2" stroke-dasharray="8,5" opacity="0.5"
/>
<line
x1={6 * scaleX} y1="0"
x2={6 * scaleX} y2={H}
stroke="darkred" stroke-width="2" stroke-dasharray="8,5" opacity="0.5"
/>

<!-- Blue team players -->
{#each blue.players as player}
<circle cx={player.cx} cy={player.cy} r="18" fill="blue" opacity="0.25" />
<text
x={player.cx} y={player.cy}
text-anchor="middle" dominant-baseline="middle"
font-size="11" font-weight="bold" fill="black"
>{player.name}</text>
{/each}

<!-- Red team players -->
{#each red.players as player}
<circle cx={player.cx} cy={player.cy} r="18" fill="red" opacity="0.25" />
<text
x={player.cx} y={player.cy}
text-anchor="middle" dominant-baseline="middle"
font-size="11" font-weight="bold" fill="black"
>{player.name}</text>
{/each}

<!-- Ball -->
<circle cx={ballX} cy={ballY} r="10" fill="yellow" stroke="darkblue" stroke-width="2" />

<!-- Blue bench -->
{#if blue.bench.length > 0}
{#each blue.bench as name, i}
{@const bx = (i / Math.max(blue.bench.length - 1, 1)) * (W * 0.35)}
<text
x={bx} y={benchY}
text-anchor="middle" dominant-baseline="middle"
font-size="11" font-weight="bold" fill="blue"
>{name}</text>
{/each}
{/if}

<!-- Red bench -->
{#if red.bench.length > 0}
{#each red.bench as name, i}
{@const bx = W - (i / Math.max(red.bench.length - 1, 1)) * (W * 0.35)}
<text
x={bx} y={benchY}
text-anchor="middle" dominant-baseline="middle"
font-size="11" font-weight="bold" fill="red"
>{name}</text>
{/each}
{/if}

<!-- Bench label -->
{#if blue.bench.length > 0 || red.bench.length > 0}
<text
x={W / 2} y={benchY}
text-anchor="middle" dominant-baseline="middle"
font-size="10" fill="#666"
>— bench —</text>
{/if}
</svg>
</div>

<style>
.court-wrapper {
display: flex;
         justify-content: center;
         margin-top: 1.5rem;
}

.court-svg {
width: 100%;
       max-width: 500px;
       border-radius: 8px;
}
</style>
