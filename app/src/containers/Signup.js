import React, { useCallback} from "react";
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

const Signup = ({ history }) => {

  const classes = useStyles();


  const handleLogin = useCallback(
    async event => {
      event.preventDefault();
      const { username, password } = event.target.elements;
      var bodyFormData = new FormData();
      bodyFormData.append('username', username.value);
      bodyFormData.append('password', password.value);
      API.post('/signup', bodyFormData)
        .then(response => {
          console.log(response.data);
          const jsonMessage = response.data.message;
          if (jsonMessage === "error") {
            alert("Username already in used! Choose another username!");
            console.log("User entered wrong username/password!")
          } else {
            alert("User account created, please proceed to login!");
            history.push('/login')
          }
        })
        .catch(error => {
          alert("Username already in used! Choose another username!");
          console.log(error);
        });
    },
    [history]
  );


  return (
    <React.Fragment>
      <CssBaseline />
      <Container component="main" maxWidth="xs">
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Signup
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
            <Button
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              type="submit"
            >
              Signup
          </Button>
          </form>
        </div>
        <Box mt={8}>
        </Box>
      </Container>
    </React.Fragment>
  );
}

export default Signup;