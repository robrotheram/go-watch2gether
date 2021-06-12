import store from '../index';

export * from './user.actions';
export * from './user.reducer';

export const GetUsername = () => {
  const { user } = store.getState();
  console.log(user);
  return user.username;
};

export const GetWatcher = () => {
  const { user } = store.getState();
  return user;
};

export const GetID = () => {
  const { user } = store.getState();
  console.log(user);
  return user.id;
};
