const API_BASE_URL = '/api'

export interface Account {
  ID: string
  Email: string
  FirstName: string
  LastName: string
  Nickname: string
  Picture: string
}

export async function fetchAccount(accountId: string): Promise<Account> {
  const response = await fetch(`${API_BASE_URL}/account/${encodeURIComponent(accountId)}`, {
    credentials: 'include',
  })

  if (!response.ok) {
    throw new Error('Error al obtener la cuenta')
  }

  return response.json()
}
