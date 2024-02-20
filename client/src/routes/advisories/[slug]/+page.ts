import type { PageLoad } from "./$types";

export const load: PageLoad = ({ params }) => {
  // TODO: Use the slug to receive the advisory the user wants to see.
  return {
    id: params.slug
  };
};
