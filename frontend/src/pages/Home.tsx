import { Box, Typography, Stack, Paper } from '@mui/material'
import { SearchTeamForm } from '../components/SearchTeamForm'
import SearchIcon from '@mui/icons-material/Search'

export function Home() {
  return (
    <Box>
      <Stack spacing={4}>
        {/* Hero Section */}
        <Paper
          elevation={0}
          sx={{
            p: 4,
            background: 'linear-gradient(135deg, #00C853 0%, #6A1B9A 100%)',
            color: 'white',
            borderRadius: 4,
          }}
        >
          <Stack spacing={2} alignItems="center">
            <SearchIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Encuentra tu Equipo
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Busca equipos por deporte y nombre para ver su información, estadísticas y miembros
            </Typography>
          </Stack>
        </Paper>

        {/* Search Form */}
        <SearchTeamForm />
      </Stack>
    </Box>
  )
}

