const INITIAL_STATE = {
  playlists:[],
  loading:false
}
export const playlistsReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case "GET_PLAYLISTS":
      return {
        ...state, playlists : action.playlists, loading:false
      };
    default: return state;
  }
};

