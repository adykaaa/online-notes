
function SignUp() {
  return (
  <div class="login-box">
  <h2>Register</h2>
  <form>
    <div class="user-box">
      <input type="text" name="firstname" required=""/>
      <label>First Name</label>
    </div>
    <div class="user-box">
      <input type="text" name="lastname" required=""/>
      <label>Last Name</label>
    </div>
    <div class="user-box">
      <input type="text" name="username" required=""/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required=""/>
      <label>Password</label>
    </div>
    <a>
      <button class="register-button">Register!</button>
    </a>
  </form>
</div>
  );
}

export default SignUp;
