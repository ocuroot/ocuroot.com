---
title: "Winding Down Ocuroot"
slug: winding-down-ocuroot
excerpt: "After 18 months of working on Ocuroot, I think it's time to move on. It's been a heck of a ride and quite the learning experience!"
coverImage:
  src: "/assets/blog/winding-down-ocuroot/cover.jpg" 
  alt: "Close up of car start button"
  credit: "Photo by Ananthu Ananthu"
  creditURL: "https://www.pexels.com/photo/car-start-button-in-close-up-photography-12551173/"
date: "2026-01-01T10:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/iac-bootstrap-problem/cover.jpg"
---

I've been working on Ocuroot since June of 2024. After 18 months of solo learning, building, marketing and re-learning, it's time to call it quits.

Ocuroot came about to solve a very personal problem I had with CI/CD, specifically challenges around scaling to many environments (see [Why Ocuroot?](https://www.ocuroot.com/blog/04-why-ocuroot/)). In retrospect, perhaps this problem was a little *too* personal. Still, I persevered at it and built something I was almost happy using day-to-day. Emphasis on the almost.

I consider myself fortunate that I was in a position to handle not getting a paycheck for a year or so. In fact, I could probably push a bit further, but I've reached a point where I don't think more time and effort would cause a breakthrough for this particular idea.

A lot of shutdown posts I've seen are described as "the hardest post I've ever had to write", and the decision itself even harder. Now that I've made the decision, I don't feel like that. I've tested a theory, and disproved it, now I can move on. Plus I've discussed this with enough people in recent weeks that I know exactly what I want to say.

# Why it didn't work out

I could point to any number of reasons why Ocuroot didn't work out, some related to the idea itself, some to my own temperament. What most of them boil down to, however, is that I kept pursuing the idea when all the signs said I wouldn't be able to scale it quickly enough.

Timing was a serious challenge. I interviewed engineers and leaders at nearly 50 companies in the early stages was that they either already had their own solution for their CI/CD scaling problems, or scaling just wasn't a problem yet. The window for having a "hair on fire" problem was probably a couple of weeks every 5 years. This would leave me either having to get very lucky with outreach, or playing a very long game of visibility.

There were also indications that the problem space itself was too broad. As I started
talking about PoCs, I discovered that every company had a slightly different need. This would leave me having to build different features for each customer, which didn't seem tenable for me building solo.

Pricing is hard, and I discovered first-hand just how hard it could be. Given that existing CI/CD solutions are either very cheap or open source, I could never find a
pricing model that I was happy with. I could see a future where I was basically doing arbitrage on compute costs. This left every pricing model I came up with either looking overpriced or unsustainable.

Finally, I'm not a great salesperson. I've been told (well, flattered) that I'm good at telling stories. I even got down a fairly decent pitch for Ocuroot that at least left people wanting to know more. But I was terrible at asking for money. Believe it or not, that last part is really important for running a business.

# Why it was worth it

While the last 18 months could have been considered a failure, they were far from a waste.

First off, I had a perfect excuse (and time) to blog more. I've been writing occasionally since early in my career, but never regularly. Since starting Ocuroot, I've posted weekly without fail. Earlier this year I started splitting my writing between my personal blog - [The Friday Deploy](https://thefridaydeploy.substack.com/)) - and the Ocuroot blog, which I cross posted to [ocuroot.com](https://www.ocuroot.com/blog/) and LinkedIn as [Scaling Deployments](https://www.linkedin.com/newsletters/scaling-deployments-7307762478157885441/). It's been great to flex my writing and research muscles, and I'm determined to keep this habit going.

I've also been motivated to attend local tech events and met a lot of people in-person. I mean a *lot*. Over 1000 based on my LinkedIn connection count. A few of these people I've seen again and again at events and conferences and would consider them good new friends. I'll continue to attend favorite events and looking for new ones, but might not maintain the pace of three a week. Not sure I have that in me.

Along with attending events, I've sought out some opportunities a little public speaking. This included presenting [Migrations for Humans](https://docs.google.com/presentation/d/1jJ_BXXF2_WvSrZ-g0K1eU2ycer6-yoKwAUFI2HG0FXI/edit?usp=sharing) at Google NY SRE Tech Talks, [When Colleagues are Customers](https://docs.google.com/presentation/d/1woorPCwwGq6vLqB4qTCpLHg4mUs5HuW_x9_bnrGN_Dc/edit?usp=sharing) at Platform Engineering New York and talking a little about an [early iteration of Ocuroot](https://www.ocuroot.com/blog/07-enabling-pipeline-visualization/) at infra.nyc. I also guested on podcasts like [Tern Stories](https://www.youtube.com/watch?v=hGx-Td8ESgM), and [Coding Chats](https://www.youtube.com/watch?v=8npL7b8AwiY).

The code I wrote also isn't going to waste. Everything is [open source](https://github.com/ocuroot/ocuroot), and I'm starting to extract some of the more interesting pieces for other projects. I'm a month into a [blog series](https://thefridaydeploy.substack.com/p/can-git-back-a-rest-api-part-1-the) on using Git for storage and have a few ideas cooking for things I could do with my Starlark helpers. 

Having more flexibility with my time has also been really nice. I've been taking most Fridays off and working Sundays to fit in better with my wife's work schedule. Getting that extra time together has been priceless. It's also been nice to be out and about in the city during weekday afternoons when things are quieter.

# What I learned

Whenever something doesn't work out, there are always lessons to be learned. And I learned my fair share from this experience. There are the obvious ones: selling is hard, building solo is lonely. But also a few less obvious ones.

I was surprised at how generous people could be with their time. It was way easier than I expected to get feedback from friends and former colleagues, and even to have them jump on a call to test out my latest quickstart.

Video content is surprisingly powerful. Early on, I had a 1 minute introduction video on the Ocuroot home page that I wasn't happy with, but when I took it down, mailing list signups basically dropped off completely.

# What I could have done differently

Do I think I could have forced Ocuroot into working? Probably not. Could someone else have done it? There's probably someone out there who would have been more successful. Probably many people if I'm totally honest with myself! I could certainly have done a few things differently, either to make Ocuroot more successful or come to my conclusions about the idea more quickly.

I mentioned earlier that there were problem signs showing during early discovery interviews. This was within the first few months and should have been a huge red flag. I didn't appreciate just how difficult finding real customers would be, I just chalked it up to "not knowing the right people" when setting up those initial conversations. If I had something I could repeatably sell, the excitement probably should have been palpable.

As I was starting to think about making changes, I picked up a few consulting engagements. Partly to bring in at least a little money for the year, and partly to force me to do something a little different so I didn't burn out. These were short engagements, a few days each at most. But even in that short time I saw real-world DevOps problems that needed real expertise to address. In retrospect, I wish I'd done some of these engagements sooner, maybe even built some bespoke solutions for these problems that could have evolved into a real product. I've seen a lot of advice floating about to "sell before building", and I've been pretty skeptical of it. But now I see that you can approach it as selling your own expertise and experience, deliver something bespoke and *then* build something more generic you can sell at scale.

I also could have spent more money. I was very careful about every penny, only spending about half of the initial cash I put in. There were definitely enough fund to try more things, maybe run some ads to test messaging and concepts, or invest in hosting to offer more powerful services earlier on.

# What's next?

As I mentioned, Ocuroot itself will stay open source for the time being. I'll be using that work to inspire a few more side projects and experiments that I'll be talking about on my personal blog.

Career wise, I've been looking at a few full time roles and might pick up a bit more consulting work in the mean time. The people I've met during this journey have been incredibly kind in connecting me with opportunities, for which I'm very grateful. 2026 is looking like an exciting year, and it wouldn't be half as exciting if I hadn't taken the leap to give Ocuroot a try.