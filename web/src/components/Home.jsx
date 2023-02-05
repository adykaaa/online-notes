import { useContext, useState } from "react";
import { useToast, Container } from '@chakra-ui/react'
import ProSidebar from "./Sidebar"
import { UserContext } from "./UserContext";
import NewNote from "./NewNote"

function Home() {

  return (
    <>
    <Container minH="100vh" minW='100vw' display="flex" margin="0 0 0 0" padding="0 0 0 0" overflow="hidden">
      <ProSidebar/>
      <NewNote display="flex"/>
    </Container>
    </>
    )
}

export default Home