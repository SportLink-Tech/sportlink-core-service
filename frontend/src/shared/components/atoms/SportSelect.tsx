import { TextField, MenuItem } from '@mui/material'
import { Sport } from '../../types/team'

/**
 * Atom: Sport Select Component
 * Reusable basic component following Atomic Design
 */
interface SportSelectProps {
  value: Sport
  onChange: (value: Sport) => void
  required?: boolean
  sports: Sport[]
}

export function SportSelect({ value, onChange, required = false, sports }: SportSelectProps) {
  return (
    <TextField
      select
      label="Deporte"
      value={value}
      onChange={(e) => onChange(e.target.value as Sport)}
      required={required}
      fullWidth
    >
      {sports.map((sport) => (
        <MenuItem key={sport} value={sport}>
          {sport}
        </MenuItem>
      ))}
    </TextField>
  )
}

