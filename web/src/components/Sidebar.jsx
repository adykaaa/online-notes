import { Sidebar, Menu, MenuItem, useProSidebar } from "react-pro-sidebar";
import AddCircleOutlinedIcon from '@mui/icons-material/AddCircleOutlined';
import EventNoteOutlinedIcon from '@mui/icons-material/EventNoteOutlined';
import MenuOutlinedIcon from "@mui/icons-material/MenuOutlined";
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import axios from "axios";
import { useContext } from "react";
import { UserContext } from "./UserContext";
import ShowToast from './Toast'
import { useToast } from '@chakra-ui/react'

function ProSidebar() {

  const toast = useToast()
  const { collapseSidebar } = useProSidebar();
  const { user, dispatch } = useContext(UserContext)
  
  const logout = () => {
    axios.post("http://localhost:8080/logout", {username: user}).then(response => {
        if (response.status == 200) {
            dispatch({type: 'LOGOUT', payload: user})
            localStorage.removeItem('user')
            ShowToast(toast,"info","Successfully logged out!")
        }
    })
    .catch(function () {
      ShowToast(toast,"error","Server error while logging out, please try again later.")
    })
  }
  
  return (
    <div id="app" style={({ height: "100vh" }, { display: "flex" })}>
      <Sidebar style={{ height: "100vh", backgroundColor:"white", maxWidth:"400px", fontSize:"20px"}}>
        <Menu>
          <MenuItem
            icon={<MenuOutlinedIcon />}
            onClick={() => {
              collapseSidebar();
            }}
            style={{ textAlign: "center", color: "linear-gradient(#141e30, #243b55)", fontWeight:"bold", marginBottom:"30px" }}
          >
            <h2>Online NoteZ</h2>
          </MenuItem>

          <MenuItem style={{"marginBottom":"10px"}} icon={<EventNoteOutlinedIcon />}>My Notes</MenuItem>
          <MenuItem style={{"marginBottom":"10px"}} icon={<AddCircleOutlinedIcon />}>Create a Note!</MenuItem>
          <MenuItem style={{"marginBottom":"10px"}} icon={<AccountCircleOutlinedIcon />}>Profile</MenuItem>
          <MenuItem style={{"marginTop":"50px"}} icon={<LogoutOutlinedIcon />} onClick={() => logout()}>Log Out</MenuItem>
        </Menu>
      </Sidebar>
    </div>
  );
}

export default ProSidebar;