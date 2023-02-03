import { Sidebar, Menu, MenuItem, useProSidebar } from "react-pro-sidebar";
import AddCircleOutlinedIcon from '@mui/icons-material/AddCircleOutlined';
import EventNoteOutlinedIcon from '@mui/icons-material/EventNoteOutlined';
import MenuOutlinedIcon from "@mui/icons-material/MenuOutlined";
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';

function ProSidebar() {
  const { collapseSidebar } = useProSidebar();

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

          <MenuItem style={{"marginBottom":"15px"}} icon={<EventNoteOutlinedIcon />}>My Notes</MenuItem>
          <MenuItem style={{"marginBottom":"15px"}} icon={<AddCircleOutlinedIcon />}>Create a Note!</MenuItem>
          <MenuItem style={{"marginBottom":"15px"}} icon={<AccountCircleOutlinedIcon />}>Profile</MenuItem>
        </Menu>
      </Sidebar>
    </div>
  );
}

export default ProSidebar;