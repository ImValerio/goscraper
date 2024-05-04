import React from "react";

const Result = ({ element }) => {
  return (
    <div class="flex flex-col justify-start bg-white border shadow-sm rounded-xl m-1">
      <div class="bg-gray-100 border-b rounded-t-xl py-3 px-4 md:py-4 md:px-5 ">
        <p class="mt-1 text-sm text-gray-500 dark:text-neutral-500">{element.url}</p>
      </div>
      <div class="p-4 md:p-5">
        {element.res.map(el=>{
        return (
            <p class="mt-2 text-gray-500 dark:text-neutral-400">
                {el}
            </p>
        )
        })}
        
     </div>
    </div>
  );
};

export default Result;
