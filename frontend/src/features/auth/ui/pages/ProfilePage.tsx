import { Avatar, Box, Paper, Typography, Stack } from '@mui/material'
import { useAuth } from '../../context/AuthContext'

export function ProfilePage() {
  const { account } = useAuth()

  const displayName = account
    ? `${account.FirstName} ${account.LastName}`.trim()
    : ''

  return (
    <Box maxWidth={480} mx="auto" mt={4}>
      <Paper elevation={0} sx={{ p: 4, borderRadius: 3 }}>
        <Stack spacing={3} alignItems="center">
          <Avatar
            src={account?.Picture}
            alt={displayName}
            sx={{ width: 96, height: 96 }}
          />
          <Box textAlign="center">
            <Typography variant="h5" fontWeight={700}>
              {displayName}
            </Typography>
            <Typography variant="body1" color="text.secondary" mt={0.5}>
              {account?.Email}
            </Typography>
          </Box>
        </Stack>
      </Paper>
    </Box>
  )
}
