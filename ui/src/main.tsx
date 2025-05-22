import React from 'react'
import ReactDOM from 'react-dom/client'
import Index from './pages/index/index'
import App from './pages/app'
import './main.css'
import { BrowserRouter, Route, Routes} from 'react-router-dom';
import { AppController } from './pages/app/controller';
import { ErrroPage } from './pages/app/error'
import PlaylistPage from './pages/app/playlist'
import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

const queryClient = new QueryClient()

const Router = () => {
  return (
      <BrowserRouter>
          <Routes>
            
            <Route path="app" element={<App />}>
              <Route path=":id" element={<AppController />}/>
              <Route path=":id/playlists" element={<PlaylistPage />}/>
              <Route index element={<ErrroPage/>}/>
            </Route>
            <Route path="/" element={<Index />} />
          </Routes>
      </BrowserRouter>
  )
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <Router />
    </QueryClientProvider>
  </React.StrictMode>,
)