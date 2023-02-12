import React from 'react';
import ReactDOM from 'react-dom/client';
import './components/login.css';
import App from './App.jsx'; 
import {UserContextProvider} from './components/UserContext'


const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <UserContextProvider>
        <App />
    </UserContextProvider>
  </React.StrictMode>
);
