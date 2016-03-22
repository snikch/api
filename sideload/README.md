# Sideload

Sideloading defines a mechanism for retrieving related entities for a collection without having to embed them in the entity themselves.

An example is a list of posts, each with an author id. To sideload the authors, you would receive a separate array of the authors referenced in the posts array. Each author only appears once.


This library attempts to define a mechanism to automate the loading of related entities. To do so, two things are required to happen. First you need to register handlers for specific entity types, and secondly annotate your types with the relationships.

## For Example

Continuing the posts <-> authors example, you might register a handler to retrieve authors (`findAuthorsByIds` would return a map of authors for the supplied author ids keyed to the id, this is **your** code).

```go
sideload.RegisterEntityHandler("authors", findAuthorsByIds)
```

You would then annotate your post type with a tag, and register it.

```go
type Post{
  AuthorID string `sideload:"authors"`
}

type Author {
  ID string
}
// By registering the type, all reflection work is done ahead of time, and is not required to be done "just in time", where it could affect performance.
sideload.RegisterType(Post{})
```

Now you can call sideload to retrieve all of the related entities for a collection.

```go
posts := []Post{
  Post{AuthorID: "author1"},
  Post{AuthorID: "author2"},
  Post{AuthorID: "author1"},
}
entities, err := sideload.Load(posts)
if err != nil {
  panic(err)
}

// You can then expect the entities map to look like so:
reflect.DeepEqual(entities, map[string]map[string]interface{}{
  "authors": map[string]interface{}{
    "author1": Author{ID: "author1"},
    "author2": Author{ID: "author2"},
  }
})
```

# TODO

- [x] Add a mechanism to manually add pre-sideloaded entities to the payload via the context.
- [ ] Ensure pre-loaded entities aren't reloaded as part of a separate 'load' event.
