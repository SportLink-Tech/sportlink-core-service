import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  AppBar,
  Toolbar,
  Typography,
  Avatar,
  IconButton,
  Menu,
  MenuItem,
  Box,
  Container,
  Stack,
  Divider,
} from '@mui/material'
import SportsIcon from '@mui/icons-material/Sports'
import PersonIcon from '@mui/icons-material/Person'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import SearchIcon from '@mui/icons-material/Search'
import LogoutIcon from '@mui/icons-material/Logout'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

  // Usuario hardcoded
  const user = {
    name: 'Jorge',
    email: 'jorge@sportlink.com',
    avatar: 'https://i.pravatar.cc/150?img=12', // Avatar aleatorio
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleNavigate = (path: string) => {
    navigate(path)
    handleMenuClose()
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      {/* AppBar */}
      <AppBar position="sticky" elevation={0} sx={{ bgcolor: 'white', borderBottom: '1px solid #e0e0e0' }}>
        <Toolbar>
          {/* Logo y Título */}
          <IconButton
            edge="start"
            onClick={() => navigate('/')}
            sx={{ mr: 2, bgcolor: 'primary.main', '&:hover': { bgcolor: 'primary.dark' } }}
          >
            <SportsIcon sx={{ color: 'white' }} />
          </IconButton>
          
          <Typography
            variant="h6"
            component="div"
            sx={{ 
              flexGrow: 1, 
              fontWeight: 700,
              background: 'linear-gradient(45deg, #00C853 30%, #6A1B9A 90%)',
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              cursor: 'pointer'
            }}
            onClick={() => navigate('/')}
          >
            SportLink
          </Typography>

          {/* Navegación */}
          <Stack direction="row" spacing={1} sx={{ mr: 2 }}>
            <IconButton
              onClick={() => navigate('/')}
              color={location.pathname === '/' ? 'primary' : 'default'}
              sx={{ 
                bgcolor: location.pathname === '/' ? 'primary.lighter' : 'transparent'
              }}
            >
              <SearchIcon />
            </IconButton>
            <IconButton
              onClick={() => navigate('/create')}
              color={location.pathname === '/create' ? 'secondary' : 'default'}
              sx={{ 
                bgcolor: location.pathname === '/create' ? 'secondary.lighter' : 'transparent'
              }}
            >
              <AddCircleIcon />
            </IconButton>
          </Stack>

          {/* Avatar y Menú de Usuario */}
          <IconButton onClick={handleMenuOpen} sx={{ p: 0 }}>
            <Avatar
              alt={user.name}
              src={user.avatar}
              sx={{ 
                width: 40, 
                height: 40,
                border: '2px solid',
                borderColor: 'primary.main'
              }}
            />
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            sx={{ mt: 1 }}
          >
            {/* User Info */}
            <Box sx={{ px: 2, py: 1.5, minWidth: 200 }}>
              <Stack direction="row" spacing={2} alignItems="center">
                <Avatar src={user.avatar} sx={{ width: 48, height: 48 }} />
                <Box>
                  <Typography variant="subtitle1" fontWeight={600}>
                    {user.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {user.email}
                  </Typography>
                </Box>
              </Stack>
            </Box>
            
            <Divider />

            {/* Menu Items */}
            <MenuItem onClick={() => handleNavigate('/')}>
              <SearchIcon sx={{ mr: 1 }} fontSize="small" />
              Buscar Equipos
            </MenuItem>
            <MenuItem onClick={() => handleNavigate('/create')}>
              <AddCircleIcon sx={{ mr: 1 }} fontSize="small" />
              Crear Equipo
            </MenuItem>
            
            <Divider />

            <MenuItem onClick={handleMenuClose}>
              <PersonIcon sx={{ mr: 1 }} fontSize="small" />
              Mi Perfil
            </MenuItem>
            <MenuItem onClick={handleMenuClose}>
              <LogoutIcon sx={{ mr: 1 }} fontSize="small" />
              Cerrar Sesión
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Contenido Principal */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          bgcolor: 'background.default',
          py: 4,
        }}
      >
        <Container maxWidth="lg">
          {children}
        </Container>
      </Box>

      {/* Footer */}
      <Box
        component="footer"
        sx={{
          py: 3,
          px: 2,
          mt: 'auto',
          bgcolor: 'white',
          borderTop: '1px solid #e0e0e0',
        }}
      >
        <Container maxWidth="lg">
          <Typography variant="body2" color="text.secondary" align="center">
            © 2025 SportLink - Plataforma de gestión deportiva
          </Typography>
        </Container>
      </Box>
    </Box>
  )
}

