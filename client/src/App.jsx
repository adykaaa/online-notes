import { Navigate, Route, Routes } from "react-router-dom";
import  Signup  from "./components/Signup.jsx";
import Home from "./pages/Home"
import Login from "./components/Login.jsx"
import { useContext} from "react";
import { BrowserRouter } from 'react-router-dom';
import { ProSidebarProvider } from "react-pro-sidebar";
import { ChakraProvider,extendTheme } from '@chakra-ui/react';
import { UserContext, UserContextProvider } from './components/UserContext'
import ViewNotes from "./components/ViewNotes.jsx";


function App() {

  const theme = extendTheme({
    styles: {
      global: () => ({
        body: {
          bg: "",
        },
      }),
    },
  });

  const {user} = useContext(UserContext)

  return (
    <ChakraProvider theme={theme}>
      <BrowserRouter>
        <ProSidebarProvider>
        <Routes>
          <Route path ="/" element={user ? <Home/> : <Navigate to="/login" />} />
          <Route path ="/notes" element={user ? <ViewNotes/> :<Navigate to="/login"/>} />
          <Route path ="/register" element={<Signup />} />
          <Route path ="/login" element={!user ? <Login /> : <Navigate to="/" />} />
        </Routes>
        </ProSidebarProvider>
      </BrowserRouter>
    </ChakraProvider>
  )
}

export default App;
