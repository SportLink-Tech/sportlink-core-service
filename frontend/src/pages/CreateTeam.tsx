import { Box, Typography, Stack, Paper } from '@mui/material'
import { CreateTeamForm } from '../components/CreateTeamForm'
import AddCircleIcon from '@mui/icons-material/AddCircle'

export function CreateTeam() {
  return (
    <Box>
      <Stack spacing={4}>
        {/* Hero Section */}
        <Paper
          elevation={0}
          sx={{
            p: 4,
            background: 'linear-gradient(135deg, #6A1B9A 0%, #00C853 100%)',
            color: 'white',
            borderRadius: 4,
          }}
        >
          <Stack spacing={2} alignItems="center">
            <AddCircleIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Crear Nuevo Equipo
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Registra un nuevo equipo deportivo con su información, categoría y miembros
            </Typography>
          </Stack>
        </Paper>

        {/* Create Form */}
        <CreateTeamForm />
      </Stack>
    </Box>
  )
}

