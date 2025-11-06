import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Auth } from './Auth/Auth';
import { FeedPage } from './FeedPage/FeedPage';
import { UserPostsPage } from './UserPostsPage/UserPostsPage';
import { useAuth } from './hooks/useAuth';

function App() {
  const { isAuthenticated, username, loading: authLoading, login, logout } = useAuth();

  if (authLoading) {
    return null;
  }

  if (!isAuthenticated || !username) {
    return <Auth onAuthSuccess={login} />;
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route 
          path="/" 
          element={<FeedPage username={username} onLogout={logout} />} 
        />
        <Route 
          path="/users/:id" 
          element={<UserPostsPage />} 
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;

