import React from "react";

const Modal = ({ data, setShowModal }) => {
  return (
    <div className="modal max-w-sm min-w-sm flex justify-between flex-col">
      <div className="w-full flex justify-end">
        <button
          className="text-xl font-bold px-2 py-1 bg-red-500 hover:bg-red-700 text-white transition-all ml-1"
          onClick={() => setShowModal(false)}
        >
          X
        </button>
      </div>
      <div className="p-5">
        <div className="text-xl">{data.urls}</div>

        <div className="text-xl text-slate-400">{data.pattern}</div>
      </div>
    </div>
  );
};

export default Modal;
