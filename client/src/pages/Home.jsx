import { Container } from '@chakra-ui/react'
import ProSidebar from "../components/Sidebar"
import NewNote from "./NewNote"

function Home() {

  return ( 
    <>
    <Container minH="100vh" minW='100vw' display="flex" margin="0 0 0 0" padding="0 0 0 0" overflow="hidden" alignItems="center">
      <ProSidebar/>
      <NewNote/>
    </Container>
    </>
    )
}

export default Home