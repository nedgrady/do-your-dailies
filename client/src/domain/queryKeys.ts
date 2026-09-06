// Every key includes the acting user id (null for the real/default user, a
// number while impersonating via the admin simulator) right after the
// resource name, so cached data never bleeds between users sharing this
// app - and so that invalidating `all(userId)` still prefix-matches every
// more specific key below it for that same user.
export const queryKeys = {
  chores: {
    all: (userId: number | null) => ['chores', userId] as const,
    list: (userId: number | null) => ['chores', userId, 'list'] as const,
    completions: (userId: number | null, start: Date, end: Date) =>
      [
        'chores',
        userId,
        'completions',
        start.toISOString(),
        end.toISOString(),
      ] as const,
    queue: (userId: number | null) => ['chores', userId, 'queue'] as const,
    insights: (userId: number | null) =>
      ['chores', userId, 'insights'] as const,
  },
}
