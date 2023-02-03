import { Navigate, Route, Routes } from "react-router-dom";
import  Signup  from "./components/Signup.jsx";
import Home from "./components/Home"
import Login from "./components/Login.jsx"
import { useContext, useState, useMemo} from "react";
import { BrowserRouter } from 'react-router-dom';
import { ProSidebarProvider } from "react-pro-sidebar";
import { ChakraProvider,extendTheme } from '@chakra-ui/react';
import { UserContext } from "./components/UserContext";

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
  const [user, setUser] = useState(null);
  const userValue = useMemo(() => ({ user, setUser }), [user, setUser]);

  console.log(userValue.user)

  return (
  <UserContext.Provider value={userValue}>
    <ChakraProvider theme={theme}>
      <BrowserRouter>
        <ProSidebarProvider>
          <Routes>
            <Route path="/" element={<Login />} />
            <Route path ="/register" element={<Signup />} />
            <Route path ="/home" element={userValue.user === "adykaaa" ? <Home /> : <Navigate to="/" />} />
          </Routes>
        </ProSidebarProvider>
      </BrowserRouter>
    </ChakraProvider>
  </UserContext.Provider>
  )
}

export default App;
