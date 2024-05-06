import React from "react";
import { useAppContext } from "../App";

const SetupItem = ({ item }) => {
  const { setUrls, setPattern, removeSetup } = useAppContext();
  return (
    <div className="flex w-100 my-1 justify-between">
      <div className="flex  grow-[2] w-100 ">
        <p className="text-xl">{item.urls}</p>
        <p className="text-xl ml-5">{item.pattern}</p>
      </div>

      <div className="flex grow justify-end">
        <button
          className="text-xl font-bold px-2 py-1 bg-green-500 hover:bg-green-700 text-white transition-all  ml-1"
          onClick={() => {
            setUrls(item.urls);
            setPattern(item.pattern);
          }}
        >
          <img src="upload.svg" className="svg-white" />
        </button>
        <button
          className="text-xl font-bold px-2 py-1 bg-red-500 hover:bg-red-700 text-white transition-all ml-1"
          onClick={() => removeSetup(item.urls, item.pattern)}
        >
          X
        </button>
      </div>
    </div>
  );
};

export default SetupItem;
