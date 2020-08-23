import React, { useCallback } from "react";
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Box from '@material-ui/core/Box';
import Language from '@material-ui/icons/Language';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';

import API from "../utils/API";

import Header from "../components/Header";


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

export default function AddItem() {

  const classes = useStyles();


  const handleLogin = useCallback(
    async event => {
      event.preventDefault();
      const { itemURL } = event.target.elements;
      console.log(itemURL.value);
      const urlArr = itemURL.value.split("/");
      console.log(urlArr[urlArr.length - 1]);
      const item_detail = urlArr[urlArr.length - 1];
      const itemDetailArr = item_detail.split(".");
      if (itemDetailArr.length != 3) {
        alert("Please enter a valid URL! \n E.g https://shopee.sg/Vans-Couple-Models-Cotton-Breathable-Short-sleeved-Short-Tee-Round-Neck-Shirt-i.42181947.3602381253");
      } else {
        const item_name = itemDetailArr[0].split("-").join(" ").slice(0, itemDetailArr[0].length - 2);
        const shop_id = itemDetailArr[1];
        const item_id = itemDetailArr[2];
        console.log(item_name);
        console.log(item_id);
        console.log(shop_id);
        var bodyFormData = new FormData();
        bodyFormData.append('item_id', item_id);
        bodyFormData.append('item_name', item_name);
        bodyFormData.append('shop_id', shop_id);
        API.post('/user/add-new-item', bodyFormData, { withCredentials: true })
          .then(response => {
            console.log(response.data);
            const jsonMessage = response.data.message;
            if (jsonMessage === "error") {
              alert("Item already added to your watchlist!");
              console.log("Item already added to your watchlist!")
            } else {
              alert("Item added to watchlist!");
            }
          })
          .catch(error => {
            alert("Please enter a valid URL! \n E.g https://shopee.sg/Vans-Couple-Models-Cotton-Breathable-Short-sleeved-Short-Tee-Round-Neck-Shirt-i.42181947.3602381253");
            console.log(error);
          });
      }


    },
  );


  return (
    <React.Fragment>
      <CssBaseline />
      <Header classes={classes} />
      <Container component="main" maxWidth="xs">
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <Language />
          </Avatar>
          <Typography component="h1" variant="h5">
            Shopee URL for Item
        </Typography>
          <form onSubmit={handleLogin} className={classes.form} noValidate>
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="itemURL"
              label="Shopee URL"
              name="itemURL"
              autoFocus
            />
            <Button
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              type="submit"
            >
              Add Item
          </Button>
          </form>
        </div>
        <Box mt={8}>
        </Box>
      </Container>
    </React.Fragment>
  );
}