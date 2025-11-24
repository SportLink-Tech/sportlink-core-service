import { Routes, Route } from 'react-router-dom'
import { Layout } from './components/Layout'
import { TeamModule } from './features/team/TeamModule'
import { SearchTeamPage } from './features/team/ui/pages/SearchTeamPage'
import { CreateTeamPage } from './features/team/ui/pages/CreateTeamPage'

/**
 * App Component
 * Wraps the application with TeamModule for Dependency Injection
 * Following Atomic Hexagonal Architecture
 */
function App() {
  return (
    <TeamModule>
      <Layout>
        <Routes>
          <Route path="/" element={<SearchTeamPage />} />
          <Route path="/create" element={<CreateTeamPage />} />
        </Routes>
      </Layout>
    </TeamModule>
  )
}

export default App

