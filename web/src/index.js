import React from 'react';
import ReactDOM from 'react-dom/client';
import './components/login.css';
import Login from './components/Login';
import App from './App.jsx'; 
import { BrowserRouter } from 'react-router-dom';
import { ProSidebarProvider } from "react-pro-sidebar";
import Context from "./components/UserContext";
import { ChakraProvider,extendTheme } from '@chakra-ui/react';

const theme = extendTheme({
  styles: {
    global: () => ({
      body: {
        bg: "",
      },
    }),
  },
});

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <ProSidebarProvider>
    <BrowserRouter>
    <ChakraProvider theme={theme}>
      <Context>
        <App />
      </Context>
    </ChakraProvider>
    </BrowserRouter>
    </ProSidebarProvider>
  </React.StrictMode>
);
