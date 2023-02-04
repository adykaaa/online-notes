import { useContext, useState } from "react";
import { useToast, Container } from '@chakra-ui/react'
import ProSidebar from "./Sidebar"
import { UserContext } from "./UserContext";

function Home() {

  return (
    <>
    <Container minH="full" minW='full' display="flex" justifyContent="space-between" margin="0 0 0 0" padding="0 0 0 0">
      <ProSidebar/>
    </Container>
    </>
    )
}

export default Home