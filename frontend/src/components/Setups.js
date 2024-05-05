import React from "react";
import SetupItem from "./SetupItem";

const Setups = ({ items }) => {
  return (
    <ul>
      {items.map((el) => (
        <SetupItem item={el} />
      ))}
    </ul>
  );
};

export default Setups;
