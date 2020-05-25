import React, { useState, useContext, useEffect } from 'react';
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
import API from "../utils/API";

import Header from '../components/Header';
import { AuthContext } from "../utils/Auth.js";

// infinitescroll component
// import InfiniteScroll from 'react-infinite-scroller';
import { useInfiniteScroll } from 'react-infinite-scroll-hook';

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


const ARRAY_SIZE = 20;

function loadItems(prevArray = [], startCursor = 0) {
  return new Promise(resolve => {
    let newArray = prevArray;
    const params = {
      offset: startCursor + ARRAY_SIZE,
      limit: 40,
    };
    API.get('/item/get-items', { params })
      .then(response => {
        console.log(response.data);
        if (response.data.message === "success") {
            const jsonData = response.data.items;
            jsonData.forEach((doc) => newArray.push(doc));
            console.log(newArray);
          }
          resolve(newArray);
      })
      .catch(error => {
        console.log(error);
      });
      
    }
  );
}

export default function ItemsList() {
  const classes = useStyles();

  // Set default as null so we
  // know if it's still loading
  // const [itemsInfo, setItems] = useState(null);
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState([]);


  const { currentUser } = useContext(AuthContext);

  // Initialize with listening to our
  // API. The second argument
  // with the empty array makes sure the
  // function only executes once
  // const listenForItemInfo = () => {
  //   const params = {
  //     offset: scrollInfo.offset,
  //     limit: 60,
  //   };
  //   API.get('/item/get-items', { params })
  //     .then(response => {
  //       console.log(response.data);
  //       if (response.data.message === "success") {
  //         if (!itemsInfo) {
  //           const allItems = [];
  //           const jsonData = response.data.items;
  //           jsonData.forEach((doc) => allItems.push(doc));
  //           console.log(allItems);
  //           setItems(allItems);
  //           const info = scrollInfo;
  //           info.offset += 60;
  //           setScroll(info);
  //         } else {
  //           const allItems = itemsInfo;
  //           const jsonData = response.data.items;
  //           jsonData.forEach((doc) => allItems.push(doc));
  //           console.log(allItems);
  //           setItems(allItems);
  //           const info = scrollInfo;
  //           info.offset += 60;
  //           setScroll(info);
  //         }
  //         console.log(scrollInfo);

  //       }
  //     })
  //     .catch(error => {
  //       console.log(error);
  //     });
  // }

  function handleLoadMore() {
    setLoading(true);
    loadItems(items, items.length).then(newArray => {
      setLoading(false);
      setItems(newArray);
    });
  }

  // useEffect(() => {
  //   listenForItemInfo();
  // }, [itemsInfo]);

  const infiniteRef = useInfiniteScroll({
    loading,
    hasNextPage: true,
    onLoadMore: handleLoadMore,
  });


  const addItem = (item_id, item_name) => {
    var bodyFormData = new FormData();
    bodyFormData.append("item_id", item_id);
    bodyFormData.append("item_name", item_name);
    API.post('/user/add-item', bodyFormData, { withCredentials: true })
      .then(response => {
        console.log(response.data);
        if (response.data.message === "success") {
          alert("Added item to your watchlist!")
        } else if (currentUser) {
          alert("The item is already in your watchlist")
        } else {
          alert("Please login/signup!")
        }
      })
      .catch(error => {
        console.log(error);
      });
  }

  // if (items.length == 0) {
  //   return (
  //     <div>
  //       Please wait...
  //     </div>
  //   )
  // }

  return (
    <React.Fragment>
      <Header classes={classes} />
      <List ref={infiniteRef} className={classes.root}>
        <Container className={classes.listContainer} maxWidth="xs">
          {items.map(({ item_name, item_price, item_id, shop_id }, index) => (
            <Card className={classes.card} key={index}>
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
                <Button size="small" color="primary" variant="outlined" className={classes.cardButton} onClick={addItem.bind(this, item_id, item_name)}>
                  Add to Watchlist
                  </Button>
              </CardActions>
            </Card>
          ))}
        </Container>
      </List>
    </React.Fragment>
  );
}
