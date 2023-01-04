import { Link } from "react-router-dom";
import './login.css';

function Login() {
  return (
  <div class="login-box">
  <h2>Login to OnlineNotes!</h2>
  <form>
    <div class="user-box">
      <input type="text" name="username" required=""/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required=""/>
      <label>Password</label>
    </div>
    <a>
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <button class="submit-button">Submit</button>
    </a>
  </form>
    <div class="signup">
      New here? <Link to="/register">Sign Up!</Link>
    </div>
</div>
  );
}

export default Login;

