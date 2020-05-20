import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardActions from '@material-ui/core/CardActions';
import Button from '@material-ui/core/Button';
import ListItemText from '@material-ui/core/ListItemText';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';

import Header from '../components/Header';

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
    backgroundColor: theme.palette.background.paper,
  },
  listContainer: {
    paddingTop: theme.spacing(8),
    paddingBottom: theme.spacing(8),
  },
  card: {
    height: '100%',
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    boxShadow: theme.shadows[2],
    marginBottom: theme.spacing(4),
    marginTop: theme.spacing(4),
  },
  cardContent: {
    flexGrow: 1,
  },
  inline: {
    display: 'inline',
  },
  cardButton: {
    horizontalAlign: "right",
  }
}));

export default function ItemsList() {
  const classes = useStyles();

  // Set default as null so we
  // know if it's still loading
  const [itemsInfo, setItems] = useState(null);

  // Initialize with listening to our
  // API. The second argument
  // with the empty array makes sure the
  // function only executes once
  useEffect(() => {
    listenForItemInfo();
  }, []);

  const listenForItemInfo = () => {
    axios.get('http://localhost/api/item/get-items')
      .then(response => {
        console.log(response.data);
        const allItems = [];
        const jsonData = response.data.items;
        jsonData.forEach((doc) => allItems.push(doc));
        setItems(allItems)
      })
      .catch(error => {
        console.log(error);
      });
  }

  if (!itemsInfo) {
    return (
      <div>
        Please wait...
      </div>
    )
  }



  return (
    <React.Fragment>
      <Header classes={classes} />
      <List className={classes.root}>
        <Container className={classes.listContainer} maxWidth="xs">
          {itemsInfo.map(({ item_name, item_price, item_id, shop_id }, index) => (
            <React.Fragment>
              <Card className={classes.card}>
                <CardContent className={classes.cardContent}>
                  <ListItem alignItems="flex-start">
                    <ListItemText
                      primary={<React.Fragment>
                        <Typography
                          component="span"
                          variant="body1"
                          className={classes.inline}
                          color="textPrimary"
                        >
                          {item_name}
                        </Typography>

                      </React.Fragment>}
                      secondary={
                        <React.Fragment>
                          <br />
                          <Typography
                            component="span"
                            variant="body2"
                            className={classes.inline}
                            color="textPrimary"
                          >
                            Current Price: ${item_price / 100000.0}
                          </Typography>

                        </React.Fragment>
                      }
                    />
                  </ListItem>
                </CardContent>
                <CardActions>
                  <Button size="small" color="primary" variant="outlined" className={classes.cardButton}>
                    Add to Watchlist
                  </Button>
                </CardActions>
              </Card>

            </React.Fragment>

          ))}
        </Container>
      </List>
    </React.Fragment>
  );
}
