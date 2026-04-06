import { Box, Paper, Typography, Stack, Alert } from '@mui/material'
import { GoogleLogin } from '@react-oauth/google'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

export function LoginPage() {
  const { googleLogin } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)

  const handleSuccess = async (credentialResponse: { credential?: string }) => {
    if (!credentialResponse.credential) {
      setError('No se recibió el token de Google')
      return
    }

    const result = await googleLogin(credentialResponse.credential)

    if (result.success) {
      navigate('/')
    } else {
      setError(result.error ?? 'Error al iniciar sesión')
    }
  }

  return (
    <Box
      display="flex"
      alignItems="center"
      justifyContent="center"
      minHeight="100vh"
      sx={{ background: 'linear-gradient(135deg, #6A1B9A 0%, #00C853 100%)' }}
    >
      <Paper elevation={6} sx={{ p: 6, borderRadius: 4, maxWidth: 400, width: '100%' }}>
        <Stack spacing={4} alignItems="center">
          <Typography variant="h4" fontWeight={700} color="primary">
            SportLink
          </Typography>
          <Typography variant="body1" color="text.secondary" align="center">
            Iniciá sesión para continuar
          </Typography>

          {error && (
            <Alert severity="error" sx={{ width: '100%' }}>
              {error}
            </Alert>
          )}

          <GoogleLogin
            onSuccess={handleSuccess}
            onError={() => setError('Error al autenticar con Google')}
            useOneTap
          />
        </Stack>
      </Paper>
    </Box>
  )
}
