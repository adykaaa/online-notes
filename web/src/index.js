import React from 'react';
import ReactDOM from 'react-dom/client';
import './components/login.css';
import Login from './components/Login';
import App from './App.jsx'; 
import { BrowserRouter } from 'react-router-dom';
import Context from "./components/UserContext";

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Context>
        <App />
      </Context>
    </BrowserRouter>
  </React.StrictMode>
);
