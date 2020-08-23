import React, { useCallback, useContext, useEffect, useState } from "react";
import { withRouter, Redirect } from "react-router";
import { Link } from 'react-router-dom';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Box from '@material-ui/core/Box';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';

import API from "../utils/API";
import { AuthContext } from "../utils/Auth.js";


const useStyles = makeStyles(theme => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.primary.main,
  },
  form: {
    width: '100%',
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

const Login = ({ history }) => {

  const classes = useStyles();

  const { setCurrentUser } = useContext(AuthContext);

  const handleLogin = useCallback(
    async event => {
      event.preventDefault();
      const { username, password } = event.target.elements;
      var bodyFormData = new FormData();
      bodyFormData.append('username', username.value);
      bodyFormData.append('password', password.value);
      API.post('/login', bodyFormData, { withCredentials: true })
        .then(response => {
          console.log(response.data);
          const jsonMessage = response.data.message;
          if (jsonMessage === "error") {
            alert("Entered wrong credentials! Please try again! ");
            console.log("User entered wrong username/password!")
          } else {
            setCurrentUser(true);
            history.push("/");
          }
        })
        .catch(error => {
          alert("Entered wrong credentials! Please try again! ");
          console.log(error);
        });
    },
    [history]
  );

  const [authInfo, setAuth] = useState(null);
  const checkForAuth = () => {
    API.get('/user/is-auth', { withCredentials: true })
      .then(response => {
        console.log(response.data);
        if (response.data.message === "success") {
          setAuth(true)
        }
      })
      .catch(error => {
        console.log(error);
      });
  }

  useEffect(() => {
    checkForAuth();
  });


  if (authInfo) {
    return <Redirect to="/" />;
  }

  return (
    <React.Fragment>
      <CssBaseline />
      <Container component="main" maxWidth="xs">
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Shopee Price Tracker
        </Typography>
          <form onSubmit={handleLogin} className={classes.form} noValidate>
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="username"
              label="Username"
              name="username"
              type="username"
              autoFocus
            />
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
            />
            {/* <FormControlLabel
            control={<Checkbox value="remember" color="primary" />}
            label="Remember me"
          /> */}
            <Button
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              type="submit"
            >
              Login
          </Button>
            <Link to={'/signup'}>
              {"Don't have an account? Sign Up"}
            </Link>
          </form>
        </div>
        <Box mt={8}>
        </Box>
      </Container>
    </React.Fragment>
  );
}

export default withRouter(Login);