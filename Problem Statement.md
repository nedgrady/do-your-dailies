# Do Your Dailies

A household chore management app built for people who want to stay on top of
their home without thinking about it.

## The problem

Most chore apps make you plan. They give you a calendar, a room-by-room view,
or a list of everything that needs doing — and leave the decision of what to do
today entirely up to you. That decision fatigue is exactly what you were trying
to avoid.

Do Your Dailies takes a different approach: tell it what needs doing and how
often, and it tells you what to do today.

## Two goals, one system

### Compliance

Every chore has a cadence — clean the bathroom mirror monthly, vacuum weekly,
wipe down the kitchen daily. The app guarantees nothing gets neglected and
nothing gets over-done. You shouldn't have to remember when you last did
something, or whether it's due yet.

### Efficiency

On any given day you have limited time and energy. The app surfaces a small,
manageable queue of the highest-priority tasks — the ones most overdue relative
to their cadence. You work through the queue, tick things off, and stop. No
decisions, no guilt about what you didn't do.

## How it works

Each chore has a frequency (e.g. every 3 days) and a log of completions. The
app calculates what's due, surfaces the top N tasks for today based on your
daily cap, and gets out of the way. A always-on screen shows you today's queue.
You tick things off as you go.

If life gets busy and things slip, a dashboard shows you the backlog. You
temporarily raise your daily cap until you're caught up, then drop it back down.

## Principles

- The queue is the interface. No planning required.
- Compliance is a guarantee, not a goal. If you trust the queue, you never need
  to check a dashboard.

## Types of user

### CAF — Capacity-First

“I have a fixed amount of effort I'm willing to spend.”

User sets a target work rate, e.g. 2 chores/day.
That capacity is treated as the constraint.
The system determines what that means for their chores.
Some slippage is acceptable if their capacity isn't enough.
The key question is: “Given what I'm willing to do, how well can I maintain my chores?”

So the important outputs are things like effective cadences, utilization, and slippage.

### COF — Compliance-First

“I will do what it takes to complete chores on their required cadences.”

Chore cadences are treated as the constraint.
The user is willing to vary their effort over time, within reason.
The system determines how much work is required to remain compliant.
Slippage is the thing they're trying to avoid.
The key question is: “How much do I need to do to keep everything on schedule?”

So the important output is required work rate/capacity.
