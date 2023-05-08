import React from 'react'
import ReactDOM from 'react-dom/client'
import { ErrroPage } from './pages/app/Error'
import Index from './pages/index'
import './main.css'
import { BrowserRouter, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { UserContextProvider } from './context/user';
import PlaylistPage from './pages/app/playlists/playtlist'
import Controller from './components/Controller'
import App from './pages/app/App'

const Router = () => {
  return (
    <BrowserRouter>
      <UserContextProvider>
        <Routes>
          <Route path="app" element={<App />}>
            <Route path=":id" element={<Controller />}/>
            <Route path=":id/playlists" element={<PlaylistPage />}/>
            <Route index element={<ErrroPage/>}/>
          </Route>
          <Route path="/" element={<Index />} />
        </Routes>
      </UserContextProvider>
    </BrowserRouter>
  )
}


ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <Router />
  </React.StrictMode>,
)
