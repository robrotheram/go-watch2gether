import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import Index from './index'
import './index.css'
import { BrowserRouter, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { UserContextProvider } from './context/user';

const Router = () => {
  return (
    <BrowserRouter>
      <UserContextProvider>
        <Routes>
          <Route path="/app" element={<App />}/>
          <Route path="/app/:id" element={<App />}/>
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
