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
  Paper,
} from '@mui/material'
import SportsIcon from '@mui/icons-material/Sports'
import PersonIcon from '@mui/icons-material/Person'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import LogoutIcon from '@mui/icons-material/Logout'
import GroupAddIcon from '@mui/icons-material/GroupAdd'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import EventIcon from '@mui/icons-material/Event'

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
              <EventIcon />
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

          {/* Avatar y Menú de Usuario - Estilo Facebook */}
          <Box
            onClick={handleMenuOpen}
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              cursor: 'pointer',
              borderRadius: 2,
              p: 0.5,
              '&:hover': {
                bgcolor: 'rgba(0, 0, 0, 0.04)',
              },
            }}
          >
            <Avatar
              alt={user.name}
              src={user.avatar}
              sx={{ 
                width: 40, 
                height: 40,
              }}
            />
            <ArrowDropDownIcon sx={{ color: 'text.secondary', fontSize: 20 }} />
          </Box>

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
            PaperProps={{
              sx: {
                mt: 1.5,
                minWidth: 360,
                maxWidth: 360,
                borderRadius: 2,
                boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
                overflow: 'hidden',
              },
            }}
            MenuListProps={{
              sx: { py: 0 },
            }}
          >
            {/* User Info Header */}
            <Paper
              elevation={0}
              sx={{
                px: 2,
                py: 2,
                bgcolor: 'background.paper',
                borderBottom: '1px solid',
                borderColor: 'divider',
              }}
            >
              <Stack direction="row" spacing={2} alignItems="center">
                <Avatar src={user.avatar} sx={{ width: 56, height: 56 }} />
                <Box sx={{ flex: 1 }}>
                  <Typography variant="subtitle1" fontWeight={600} fontSize="1.05rem">
                    {user.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" fontSize="0.875rem">
                    {user.email}
                  </Typography>
                </Box>
              </Stack>
            </Paper>

            {/* Menu Items */}
            <Box sx={{ py: 0.5 }}>
              <MenuItem 
                onClick={() => handleNavigate('/create')}
                sx={{ 
                  py: 1.5,
                  px: 2,
                  '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.04)' }
                }}
              >
                <AddCircleIcon sx={{ mr: 2, fontSize: 24, color: 'text.secondary' }} />
                <Typography variant="body1">Publicar Partido</Typography>
              </MenuItem>
              <MenuItem 
                onClick={() => handleNavigate('/create-team')}
                sx={{ 
                  py: 1.5,
                  px: 2,
                  '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.04)' }
                }}
              >
                <GroupAddIcon sx={{ mr: 2, fontSize: 24, color: 'text.secondary' }} />
                <Typography variant="body1">Agregar equipo</Typography>
              </MenuItem>
            </Box>

            <Divider />

            {/* Settings & Logout */}
            <Box sx={{ py: 0.5 }}>
              <MenuItem 
                onClick={handleMenuClose}
                sx={{ 
                  py: 1.5,
                  px: 2,
                  '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.04)' }
                }}
              >
                <PersonIcon sx={{ mr: 2, fontSize: 24, color: 'text.secondary' }} />
                <Typography variant="body1">Mi Perfil</Typography>
              </MenuItem>
              <MenuItem 
                onClick={handleMenuClose}
                sx={{ 
                  py: 1.5,
                  px: 2,
                  '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.04)' }
                }}
              >
                <LogoutIcon sx={{ mr: 2, fontSize: 24, color: 'text.secondary' }} />
                <Typography variant="body1">Cerrar Sesión</Typography>
              </MenuItem>
            </Box>
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

