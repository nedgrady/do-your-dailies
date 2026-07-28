export const queryKeys = {
  chores: {
    all: () => ['chores'] as const,
    list: () => ['chores', 'list'] as const,
    completions: (start: Date, end: Date) =>
      [
        'chores',
        'completions',
        start.toISOString(),
        end.toISOString(),
      ] as const,
    queue: () => ['chores', 'queue'] as const,
    insights: () => ['chores', 'insights'] as const,
  },
}
