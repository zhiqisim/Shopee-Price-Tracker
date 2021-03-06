import React, { useState } from 'react';
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
import { Link } from 'react-router-dom';
import API from "../utils/API";

// infinitescroll component
import { useInfiniteScroll } from 'react-infinite-scroll-hook';
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
  },
  link: {
    margin: theme.spacing(1, 1.5),
    textDecoration: "none"
  },
}));

// const ARRAY_SIZE = 20;

function loadItems(prevArray = [], startCursor = 0) {
  return new Promise(resolve => {
    let newArray = prevArray;
    const params = {
      offset: startCursor,
      limit: 20,
    };
    API.get('/user/watchlist', { withCredentials: true, params})
      .then(response => {
        console.log(response.data);
        if (response.data.message === "success") {
            const jsonData = response.data.item_details;
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

export default function WatchList() {
  const classes = useStyles();

  // Set default as null so we
  // know if it's still loading
  // const [watchListInfo, setWatchList] = useState(null);
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState([]);

  // Initialize with listening to our
  // API. The second argument
  // with the empty array makes sure the
  // function only executes once
  // useEffect(() => {
  //   listenForWatchListInfo();
  // }, []);

  // const listenForWatchListInfo = () => {
  //   const params = {
  //     offset: 0,
  //     limit: 20,
  //   };
  //   API.get('/user/watchlist', { withCredentials: true, params})
  //     .then(response => {
  //       console.log(response.data);
  //       if (response.data.message === "success") {
  //         const allItems = [];
  //         const jsonData = response.data.item_details;
  //         if (jsonData != null){
  //           jsonData.forEach((doc) => allItems.push(doc));
  //         }
  //         setWatchList(allItems)
  //       }
  //     })
  //     .catch(error => {
  //       console.log(error);
  //     });
  // }

  // if (!watchListInfo) {
  //   return (
  //     <div>
  //       Please wait...
  //     </div>
  //   )
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



  return (
    <React.Fragment>
      <Header classes={classes} />
      <List ref={infiniteRef} className={classes.root}>
        <Container className={classes.listContainer} maxWidth="xs">
          {items.map(({ item_name, item_id }, index) => (
            <React.Fragment key={index}>
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
                    />
                  </ListItem>
                </CardContent>
                <CardActions>
                  <Link to={'/price/' + item_id} className={classes.link}>
                    <Button size="small" color="primary" variant="outlined" className={classes.cardButton}>
                      View Price Changelog
                  </Button>
                  </Link>
                </CardActions>
              </Card>

            </React.Fragment>

          ))}
        </Container>
      </List>
    </React.Fragment>
  );
}
