//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config'

export default [
  ...tanstackConfig,
  {
    rules: {
      'import/no-cycle': 'off',
      'import/order': 'off',
      'sort-imports': 'off',
      '@typescript-eslint/array-type': 'off',
      '@typescript-eslint/require-await': 'off',
      'pnpm/json-enforce-catalog': 'off',
      // Dialogs and sheets must go through ResponsiveDialog, which keeps the
      // title and actions on screen. Raw MUI dialogs scroll the actions away.
      'no-restricted-imports': [
        'error',
        {
          paths: [
            {
              name: '@mui/material',
              importNames: [
                'Dialog',
                'DialogTitle',
                'DialogContent',
                'DialogActions',
                'Drawer',
                'SwipeableDrawer',
                'Modal',
              ],
              message:
                'Use ResponsiveDialog (src/components/ResponsiveDialog.tsx).',
            },
          ],
          patterns: [
            {
              group: [
                '@mui/material/Dialog*',
                '@mui/material/Drawer',
                '@mui/material/SwipeableDrawer',
                '@mui/material/Modal',
              ],
              message:
                'Use ResponsiveDialog (src/components/ResponsiveDialog.tsx).',
            },
          ],
        },
      ],
    },
  },
  {
    // The one place allowed to build a raw MUI dialog.
    files: ['src/components/ResponsiveDialog.tsx'],
    rules: {
      'no-restricted-imports': 'off',
    },
  },
  {
    ignores: ['eslint.config.js', 'prettier.config.js', 'src/generated'],
  },
]
