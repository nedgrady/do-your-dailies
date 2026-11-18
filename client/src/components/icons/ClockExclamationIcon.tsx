import { SvgIcon, type SvgIconProps } from '@mui/material'

// Clock with an exclamation mark, drawn from Lucide's clock-exclamation-point.
// Stroke-based, so fill is off and stroke follows the current colour.
export default function ClockExclamationIcon(
  props: Omit<SvgIconProps, 'sx'>,
) {
  return (
    <SvgIcon
      viewBox="0 0 24 24"
      {...props}
      sx={{
        fill: 'none',
        stroke: 'currentColor',
        strokeWidth: 2,
        strokeLinecap: 'round',
        strokeLinejoin: 'round',
      }}
    >
      <path d="M12 6v6l4 2" />
      <path d="M20 12v5" />
      <path d="M20 21h.01" />
      <path d="M21.25 8.2A10 10 0 1 0 16 21.16" />
    </SvgIcon>
  )
}
