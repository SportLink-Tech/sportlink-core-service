import { TextField, MenuItem } from '@mui/material'

/**
 * Atom: Category Select Component
 * Reusable basic component following Atomic Design
 */
const CATEGORIES = [
  { value: 0, label: 'Unranked' },
  { value: 1, label: 'L1 - Principiante' },
  { value: 2, label: 'L2' },
  { value: 3, label: 'L3' },
  { value: 4, label: 'L4' },
  { value: 5, label: 'L5' },
  { value: 6, label: 'L6' },
  { value: 7, label: 'L7 - Avanzado' },
]

interface CategorySelectProps {
  value: number
  onChange: (value: number) => void
  fullWidth?: boolean
}

export function CategorySelect({ value, onChange, fullWidth = true }: CategorySelectProps) {
  return (
    <TextField
      select
      label="Categoría"
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
      fullWidth={fullWidth}
    >
      {CATEGORIES.map((cat) => (
        <MenuItem key={cat.value} value={cat.value}>
          {cat.label}
        </MenuItem>
      ))}
    </TextField>
  )
}

