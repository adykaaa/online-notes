import './login.css';

function Login() {
  return (
  <div class="login-box">
  <h2>Login to MovieStore!</h2>
  <form>
    <div class="user-box">
      <input type="text" name="username" required=""/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required=""/>
      <label>Password</label>
    </div>
    <a href="#">
      Submit
    </a>
  </form>
</div>
  );
}

export default Login;

