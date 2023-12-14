import React, { useEffect } from 'react'
import ReactDOM from 'react-dom/client'
import { ErrroPage } from './pages/app/Error'
import Index from './pages/index'
import './main.css'
import { BrowserRouter, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { UserContextProvider } from './context/user';
import PlaylistPage from './pages/app/playlists/playtlist'
import { AppController } from './components/Controller'
// import App from './pages/app/App'

const Router = () => {
  return (
    <BrowserRouter>
      <UserContextProvider>
        <Routes>
          <Route path="app" element={<App />}>
            <Route path=":id" element={<AppController />}/>
            <Route path=":id/playlists" element={<PlaylistPage />}/>
            <Route index element={<ErrroPage/>}/>
          </Route>
          <Route path="/" element={<Index />} />
        </Routes>
      </UserContextProvider>
    </BrowserRouter>
  )
}

const App = () => {
  useEffect(() => {
    const sse = new EventSource('http://localhost:8000/sse',
      { withCredentials: true });
    function getRealtimeData(data) {
      console.log(data);
    }
    sse.onmessage = e => getRealtimeData(JSON.parse(e.data));
    sse.onerror = () => {
      // error log here 
      
      sse.close();
    }
    return () => {
      sse.close();
    };
  }, []);


}


ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
