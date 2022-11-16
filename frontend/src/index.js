import React from 'react';
import ReactDOM from 'react-dom/client';
import './components/login.css';
import Login from './components/Login';
import App from './App.jsx'; 
import { BrowserRouter } from 'react-router-dom';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
