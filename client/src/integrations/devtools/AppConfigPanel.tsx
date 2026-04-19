const config = [
  { label: 'API Base URL', value: import.meta.env.VITE_API_BASE_URL },
  { label: 'Mode', value: import.meta.env.MODE },
  { label: 'MSW enabled', value: import.meta.env.VITE_ENABLE_MSW ?? 'false' },
] as const

export default function AppConfigPanel() {
  return (
    <table
      style={{
        borderCollapse: 'collapse',
        width: '100%',
        fontFamily: 'monospace',
        fontSize: 13,
      }}
    >
      <thead>
        <tr>
          <th
            style={{
              textAlign: 'left',
              padding: '6px 12px',
              borderBottom: '1px solid #444',
            }}
          >
            Key
          </th>
          <th
            style={{
              textAlign: 'left',
              padding: '6px 12px',
              borderBottom: '1px solid #444',
            }}
          >
            Value
          </th>
        </tr>
      </thead>
      <tbody>
        {config.map(({ label, value }) => (
          <tr key={label}>
            <td style={{ padding: '5px 12px', color: '#aaa' }}>{label}</td>
            <td style={{ padding: '5px 12px' }}>
              {value ?? <em style={{ color: '#888' }}>not set</em>}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
