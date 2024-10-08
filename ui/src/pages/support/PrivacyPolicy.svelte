<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { onMount } from 'svelte';

  onMount(() => window.scrollTo(0, 0));
</script>

<svelte:head>
  <title>{$LL.privacyPolicy()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="dark:text-gray-200">
    <h1 class="py-2 text-2xl">{$LL.privacyPolicy()}</h1>
    <p class="py-4 italic text-gray-500 dark:text-gray-400">
      Updated at 2024-03-27
    </p>

    <p class="py-4">
      Thunderdome (“we,” “our,” or “us”) is committed to protecting your
      privacy. This Privacy Policy explains how your personal information is
      collected, used, and disclosed by Thunderdome.
    </p>
    <p class="py-4">
      This Privacy Policy applies to our website, and its associated subdomains
      (collectively, our “Service”) alongside our application, Thunderdome. By
      accessing or using our Service, you signify that you have read,
      understood, and agree to our collection, storage, use, and disclosure of
      your personal information as described in this Privacy Policy and our
      Terms of Service.
    </p>

    <h2 class="py-2 text-xl">Definitions and key terms</h2>
    <div class="mb-4">
      <h4 class="text-lg py-2">
        To help explain things as clearly as possible in this Privacy Policy,
        every time any of these terms are referenced, are strictly defined as:
      </h4>
      <ul class="space-y-1 list-disc list-inside">
        <li>
          Cookie: small amount of data generated by a website and saved by your
          web browser. It is used to identify your browser, provide analytics,
          remember information about you such as your language preference or
          login information.
        </li>
        <li>
          Company: when this policy mentions “Company,” “we,” “us,” or “our,” it
          refers to Thunderdome, that is responsible for your information under
          this Privacy Policy.
        </li>
        <li>
          Country: where Thunderdome or the owners/founders of Thunderdome are
          based, in this case is United States
        </li>
        <li>
          Customer: refers to the company, organization or person that signs up
          to use the Thunderdome Service to manage the relationships with your
          consumers or service users.
        </li>
        <li>
          Device: any internet connected device such as a phone, tablet,
          computer or any other device that can be used to visit Thunderdome and
          use the services.
        </li>
        <li>
          IP address: Every device connected to the Internet is assigned a
          number known as an Internet protocol (IP) address. These numbers are
          usually assigned in geographic blocks. An IP address can often be used
          to identify the location from which a device is connecting to the
          Internet.
        </li>
        <li>
          Personnel: refers to those individuals who are employed by Thunderdome
          or are under contract to perform a service on behalf of one of the
          parties.
        </li>
        <li>
          Personal Data: any information that directly, indirectly, or in
          connection with other information — including a personal
          identification number — allows for the identification or
          identifiability of a natural person.
        </li>
        <li>
          Service: refers to the service provided by Thunderdome as described in
          the relative terms (if available) and on this platform.
        </li>
        <li>
          Third-party service: refers to advertisers, contest sponsors,
          promotional and marketing partners, and others who provide our content
          or whose products or services we think may interest you.
        </li>
        <li>
          Website: Thunderdome's site, which can be accessed via this URL:
        </li>
        <li>
          You: a person or entity that is registered with Thunderdome to use the
          Services.
        </li>
      </ul>
    </div>

    <h2 class="py-2 text-xl">What Information Do We Collect?</h2>
    <p class="py-4">
      We collect information from you when you visit our app, register on our
      site, place an order, subscribe to our newsletter, respond to a survey or
      fill out a form.
    </p>

    <ul class="space-y-1 list-disc list-inside">
      <li>Name / Username</li>
      <li>Email Addresses</li>
      <li>Password</li>
    </ul>
    <ul class="space-y-1 list-disc list-inside"></ul>

    <h2 class="py-2 text-xl">How Do We Use The Information We Collect?</h2>
    <div class="mb-4">
      <h4 class="text-lg py-2">
        Any of the information we collect from you may be used in one of the
        following ways:
      </h4>
      <ul class="space-y-1 list-disc list-inside">
        <li>
          To personalize your experience (your information helps us to better
          respond to your individual needs)
        </li>
        <li>
          To improve our app (we continually strive to improve our app offerings
          based on the information and feedback we receive from you)
        </li>
        <li>
          To improve customer service (your information helps us to more
          effectively respond to your customer service requests and support
          needs)
        </li>
        <li>To process transactions</li>
        <li>
          To administer a contest, promotion, survey or other site feature
        </li>
        <li>To send periodic emails</li>
      </ul>
    </div>

    <h2 class="py-2 text-xl">
      When does Thunderdome use end user information from third parties?
    </h2>
    <p class="py-4">
      Thunderdome will collect End User Data necessary to provide the
      Thunderdome services to our customers.
    </p>
    <p class="py-4">
      End users may voluntarily provide us with information they have made
      available on social media websites. If you provide us with any such
      information, we may collect publicly available information from the social
      media websites you have indicated. You can control how much of your
      information social media websites make public by visiting these websites
      and changing your privacy settings.
    </p>

    <h2 class="py-2 text-xl">
      When does Thunderdome use customer information from third parties?
    </h2>
    <p class="py-4">
      We receive some information from the third parties when you contact us.
      For example, when you submit your email address to us to show interest in
      becoming a Thunderdome customer, we receive information from a third party
      that provides automated fraud detection services to Thunderdome. We also
      occasionally collect information that is made publicly available on social
      media websites. You can control how much of your information social media
      websites make public by visiting these websites and changing your privacy
      settings.
    </p>

    <h2 class="py-2 text-xl">
      Do we share the information we collect with third parties?
    </h2>
    <p class="py-4">
      We may share the information that we collect, both personal and
      non-personal, with third parties such as advertisers, contest sponsors,
      promotional and marketing partners, and others who provide our content or
      whose products or services we think may interest you. We may also share it
      with our current and future affiliated companies and business partners,
      and if we are involved in a merger, asset sale or other business
      reorganization, we may also share or transfer your personal and
      non-personal information to our successors-in-interest.
    </p>
    <p class="py-4">
      We may engage trusted third party service providers to perform functions
      and provide services to us, such as hosting and maintaining our servers
      and the app, database storage and management, e-mail management, storage
      marketing, credit card processing, customer service and fulfilling orders
      for products and services you may purchase through the app. We will likely
      share your personal information, and possibly some non-personal
      information, with these third parties to enable them to perform these
      services for us and for you.
    </p>
    <p class="py-4">
      We may share portions of our log file data, including IP addresses, for
      analytics purposes with third parties such as web analytics partners,
      application developers, and ad networks. If your IP address is shared, it
      may be used to estimate general location and other technographics such as
      connection speed, whether you have visited the app in a shared location,
      and type of the device used to visit the app. They may aggregate
      information about our advertising and what you see on the app and then
      provide auditing, research and reporting for us and our advertisers.   We
      may also disclose personal and non-personal information about you to
      government or law enforcement officials or private parties as we, in our
      sole discretion, believe necessary or appropriate in order to respond to
      claims, legal process (including subpoenas), to protect our rights and
      interests or those of a third party, the safety of the public or any
      person, to prevent or stop any illegal, unethical, or legally actionable
      activity, or to otherwise comply with applicable court orders, laws, rules
      and regulations.
    </p>

    <h2 class="py-2 text-xl">
      Where and when is information collected from customers and end users?
    </h2>
    <p class="py-4">
      Thunderdome will collect personal information that you submit to us. We
      may also receive personal information about you from third parties as
      described above.
    </p>

    <h2 class="py-2 text-xl">How Do We Use Your Email Address?</h2>
    <p class="py-4">
      By submitting your email address on this app, you agree to receive emails
      from us. You can cancel your participation in any of these email lists at
      any time by clicking on the opt-out link or other unsubscribe option that
      is included in the respective email. We only send emails to people who
      have authorized us to contact them, either directly, or through a third
      party. We do not send unsolicited commercial emails, because we hate spam
      as much as you do. By submitting your email address, you also agree to
      allow us to use your email address for customer audience targeting on
      sites like Facebook, where we display custom advertising to specific
      people who have opted-in to receive communications from us. Email
      addresses submitted only through the order processing page will be used
      for the sole purpose of sending you information and updates pertaining to
      your order. If, however, you have provided the same email to us through
      another method, we may use it for any of the purposes stated in this
      Policy. Note: If at any time you would like to unsubscribe from receiving
      future emails, we include detailed unsubscribe instructions at the bottom
      of each email.
    </p>

    <h2 class="py-2 text-xl">How Long Do We Keep Your Information?</h2>
    <p class="py-4">
      We keep your information only so long as we need it to provide Thunderdome
      to you and fulfill the purposes described in this policy. This is also the
      case for anyone that we share your information with and who carries out
      services on our behalf. When we no longer need to use your information and
      there is no need for us to keep it to comply with our legal or regulatory
      obligations, we’ll either remove it from our systems or depersonalize it
      so that we can't identify you.
    </p>

    <h2 class="py-2 text-xl">How Do We Protect Your Information?</h2>
    <p class="py-4">
      We implement a variety of security measures to maintain the safety of your
      personal information when you place an order or enter, submit, or access
      your personal information. We offer the use of a secure server. All
      supplied sensitive/credit information is transmitted via Secure Socket
      Layer (SSL) technology and then encrypted into our Payment gateway
      providers database only to be accessible by those authorized with special
      access rights to such systems, and are required to keep the information
      confidential. After a transaction, your private information (credit cards,
      social security numbers, financials, etc.) is never kept on file. We
      cannot, however, ensure or warrant the absolute security of any
      information you transmit to Thunderdome or guarantee that your information
      on the Service may not be accessed, disclosed, altered, or destroyed by a
      breach of any of our physical, technical, or managerial safeguards.
    </p>

    <h2 class="py-2 text-xl">
      Could my information be transferred to other countries?
    </h2>
    <p class="py-4">
      Thunderdome is incorporated in United States. Information collected via
      our website, through direct interactions with you, or from use of our help
      services may be transferred from time to time to our offices or personnel,
      or to third parties, located throughout the world, and may be viewed and
      hosted anywhere in the world, including countries that may not have laws
      of general applicability regulating the use and transfer of such data. To
      the fullest extent allowed by applicable law, by using any of the above,
      you voluntarily consent to the trans-border transfer and hosting of such
      information.
    </p>

    <h2 class="py-2 text-xl">
      Is the information collected through the Thunderdome Service secure?
    </h2>
    <p class="py-4">
      We take precautions to protect the security of your information. We have
      physical, electronic, and managerial procedures to help safeguard, prevent
      unauthorized access, maintain data security, and correctly use your
      information. However, neither people nor security systems are foolproof,
      including encryption systems. In addition, people can commit intentional
      crimes, make mistakes or fail to follow policies. Therefore, while we use
      reasonable efforts to protect your personal information, we cannot
      guarantee its absolute security. If applicable law imposes any
      non-disclaimable duty to protect your personal information, you agree that
      intentional misconduct will be the standards used to measure our
      compliance with that duty.
    </p>

    <h2 class="py-2 text-xl">Can I update or correct my information?</h2>
    <p class="py-4">
      The rights you have to request updates or corrections to the information
      Thunderdome collects depend on your relationship with Thunderdome.
      Personnel may update or correct their information as detailed in our
      internal company employment policies.
    </p>
    <p class="py-4">
      Customers have the right to request the restriction of certain uses and
      disclosures of personally identifiable information as follows. You
      can contact us in order to (1) update or correct your personally
      identifiable information, (2) change your preferences with respect to
      communications and other information you receive from us, or (3) delete
      the personally identifiable information maintained about you on our
      systems (subject to the following paragraph), by cancelling your account.
      Such updates, corrections, changes and deletions will have no effect on
      other information that we maintain, or information that we have provided
      to third parties in accordance with this Privacy Policy prior to such
      update, correction, change or deletion. To protect your privacy and
      security, we may take reasonable steps (such as requesting a unique
      password) to verify your identity before granting you profile access or
      making corrections. You are responsible for maintaining the secrecy of
      your unique password and account information at all times.
    </p>
    <p class="py-4">
      You should be aware that it is not technologically possible to remove each
      and every record of the information you have provided to us from our
      system. The need to back up our systems to protect information from
      inadvertent loss means that a copy of your information may exist in a
      non-erasable form that will be difficult or impossible for us to locate.
      Promptly after receiving your request, all personal information stored in
      databases we actively use, and other readily searchable media will be
      updated, corrected, changed or deleted, as appropriate, as soon as and to
      the extent reasonably and technically practicable.
    </p>
    <p class="py-4">
      If you are an end user and wish to update, delete, or receive any
      information we have about you, you may do so by contacting the
      organization of which you are a customer.
    </p>
    <h2 class="py-2 text-xl">Sale of Business</h2>
    <p class="py-4">
      We reserve the right to transfer information to a third party in the event
      of a sale, merger or other transfer of all or substantially all of the
      assets of Thunderdome or any of its Corporate Affiliates (as defined
      herein), or that portion of Thunderdome or any of its Corporate Affiliates
      to which the Service relates, or in the event that we discontinue our
      business or file a petition or have filed against us a petition in
      bankruptcy, reorganization or similar proceeding, provided that the third
      party agrees to adhere to the terms of this Privacy Policy.
    </p>

    <h2 class="py-2 text-xl">Affiliates</h2>
    <p class="py-4">
      We may disclose information (including personal information) about you to
      our Corporate Affiliates. For purposes of this Privacy Policy, "Corporate
      Affiliate" means any person or entity which directly or indirectly
      controls, is controlled by or is under common control with Thunderdome,
      whether by ownership or otherwise. Any information relating to you that we
      provide to our Corporate Affiliates will be treated by those Corporate
      Affiliates in accordance with the terms of this Privacy Policy.
    </p>

    <h2 class="py-2 text-xl">Governing Law</h2>
    <p class="py-4">
      This Privacy Policy is governed by the laws of United States without
      regard to its conflict of laws provision. You consent to the exclusive
      jurisdiction of the courts in connection with any action or dispute
      arising between the parties under or in connection with this Privacy
      Policy except for those individuals who may have rights to make claims
      under Privacy Shield, or the Swiss-US framework.
    </p>
    <p class="py-4">
      The laws of United States, excluding its conflicts of law rules, shall
      govern this Agreement and your use of the app. Your use of the app may
      also be subject to other local, state, national, or international laws.
    </p>
    <p class="py-4">
      By using Thunderdome or contacting us directly, you signify your
      acceptance of this Privacy Policy. If you do not agree to this Privacy
      Policy, you should not engage with our website, or use our services.
      Continued use of the website, direct engagement with us, or following the
      posting of changes to this Privacy Policy that do not significantly affect
      the use or disclosure of your personal information will mean that you
      accept those changes.
    </p>

    <h2 class="py-2 text-xl">Your Consent</h2>
    <p class="py-4">
      We've updated our Privacy Policy to provide you with complete transparency
      into what is being set when you visit our site and how it's being used. By
      using our app, registering an account, or making a purchase, you hereby
      consent to our Privacy Policy and agree to its terms.
    </p>

    <h2 class="py-2 text-xl">Links to Other Websites</h2>
    <p class="py-4">
      This Privacy Policy applies only to the Services. The Services may contain
      links to other websites not operated or controlled by Thunderdome. We are
      not responsible for the content, accuracy or opinions expressed in such
      websites, and such websites are not investigated, monitored or checked for
      accuracy or completeness by us. Please remember that when you use a link
      to go from the Services to another website, our Privacy Policy is no
      longer in effect. Your browsing and interaction on any other website,
      including those that have a link on our platform, is subject to that
      website’s own rules and policies. Such third parties may use their own
      cookies or other methods to collect information about you.
    </p>
    <h2 class="py-2 text-xl">Cookies</h2>
    <p class="py-4">
      Thunderdome uses "Cookies" to identify the areas of our website that you
      have visited. A Cookie is a small piece of data stored on your computer or
      mobile device by your web browser. We use Cookies to enhance the
      performance and functionality of our app but are non-essential to their
      use. However, without these cookies, certain functionality like videos may
      become unavailable or you would be required to enter your login details
      every time you visit the app as we would not be able to remember that you
      had logged in previously. Most web browsers can be set to disable the use
      of Cookies. However, if you disable Cookies, you may not be able to access
      functionality on our website correctly or at all. We never place
      Personally Identifiable Information in Cookies.
    </p>

    <h2 class="py-2 text-xl">
      Blocking and disabling cookies and similar technologies
    </h2>
    <p class="py-4">
      Wherever you're located you may also set your browser to block cookies and
      similar technologies, but this action may block our essential cookies and
      prevent our website from functioning properly, and you may not be able to
      fully utilize all of its features and services. You should also be aware
      that you may also lose some saved information (e.g. saved login details,
      site preferences) if you block cookies on your browser. Different browsers
      make different controls available to you. Disabling a cookie or category
      of cookie does not delete the cookie from your browser, you will need to
      do this yourself from within your browser, you should visit your browser's
      help menu for more information.
    </p>
    <h2 class="py-2 text-xl">Payment Details</h2>
    <p class="py-4">
      In respect to any credit card or other payment processing details you have
      provided us, we commit that this confidential information will be stored
      in the most secure manner possible.
    </p>
    <h2 class="py-2 text-xl">Kids' Privacy</h2>
    <p class="py-4">
      We do not address anyone under the age of 13. We do not knowingly collect
      personally identifiable information from anyone under the age of 13. If
      You are a parent or guardian and You are aware that Your child has
      provided Us with Personal Data, please contact Us. If We become aware that
      We have collected Personal Data from anyone under the age of 13 without
      verification of parental consent, We take steps to remove that information
      from Our servers.
    </p>
    <h2 class="py-2 text-xl">Changes To Our Privacy Policy</h2>
    <p class="py-4">
      We may change our Service and policies, and we may need to make changes to
      this Privacy Policy so that they accurately reflect our Service and
      policies. Unless otherwise required by law, we will notify you (for
      example, through our Service) before we make changes to this Privacy
      Policy and give you an opportunity to review them before they go into
      effect. Then, if you continue to use the Service, you will be bound by the
      updated Privacy Policy. If you do not want to agree to this or any updated
      Privacy Policy, you can delete your account.
    </p>

    <h2 class="py-2 text-xl">Third-Party Services</h2>
    <p class="py-4">
      We may display, include or make available third-party content (including
      data, information, applications and other products services) or provide
      links to third-party websites or services ("Third- Party Services").
    </p>
    <p class="py-4">
      You acknowledge and agree that Thunderdome shall not be responsible for
      any Third-Party Services, including their accuracy, completeness,
      timeliness, validity, copyright compliance, legality, decency, quality or
      any other aspect thereof. Thunderdome does not assume and shall not have
      any liability or responsibility to you or any other person or entity for
      any Third-Party Services.
    </p>
    <p class="py-4">
      Third-Party Services and links thereto are provided solely as a
      convenience to you and you access and use them entirely at your own risk
      and subject to such third parties' terms and conditions.
    </p>
    <h2 class="py-2 text-xl">Tracking Technologies</h2>
    <ul class="space-y-1 list-disc list-inside">
      <li>
        Cookies
        <br /><br />
        <p class="py-4">
          We use Cookies to enhance the performance and functionality of our
          platform but are non-essential to their use. However, without these
          cookies, certain functionality like videos may become unavailable or
          you would be required to enter your login details every time you visit
          the platform as we would not be able to remember that you had logged
          in previously.
        </p>
      </li>
      <br />
    </ul>
    <h2 class="py-2 text-xl">
      Information about General Data Protection Regulation (GDPR)
    </h2>
    <p class="py-4">
      We may be collecting and using information from you if you are from the
      European Economic Area (EEA), and in this section of our Privacy Policy we
      are going to explain exactly how and why is this data collected, and how
      we maintain this data under protection from being replicated or used in
      the wrong way.
    </p>

    <h2 class="py-2 text-xl">What is GDPR?</h2>
    <p class="py-4">
      GDPR is an EU-wide privacy and data protection law that regulates how EU
      residents' data is protected by companies and enhances the control the EU
      residents have, over their personal data.
    </p>
    <p class="py-4">
      The GDPR is relevant to any globally operating company and not just the
      EU-based businesses and EU residents. Our customers’ data is important
      irrespective of where they are located, which is why we have implemented
      GDPR controls as our baseline standard for all our operations worldwide.
    </p>

    <h2 class="py-2 text-xl">What is personal data?</h2>
    <p class="py-4">
      Any data that relates to an identifiable or identified individual. GDPR
      covers a broad spectrum of information that could be used on its own, or
      in combination with other pieces of information, to identify a person.
      Personal data extends beyond a person’s name or email address. Some
      examples include financial information, political opinions, genetic data,
      biometric data, IP addresses, physical address, sexual orientation, and
      ethnicity.
    </p>
    <p class="py-4">
      The Data Protection Principles include requirements such as:
    </p>
    <ul class="space-y-1 list-disc list-inside">
      <li>
        Personal data collected must be processed in a fair, legal, and
        transparent way and should only be used in a way that a person would
        reasonably expect.
      </li>
      <li>
        Personal data should only be collected to fulfil a specific purpose and
        it should only be used for that purpose. Organizations must specify why
        they need the personal data when they collect it.
      </li>
      <li>
        Personal data should be held no longer than necessary to fulfil its
        purpose.
      </li>
      <li>
        People covered by the GDPR have the right to access their own personal
        data. They can also request a copy of their data, and that their data be
        updated, deleted, restricted, or moved to another organization.
      </li>
    </ul>

    <h2 class="py-2 text-xl">Why is GDPR important?</h2>
    <p class="py-4">
      GDPR adds some new requirements regarding how companies should protect
      individuals' personal data that they collect and process. It also raises
      the stakes for compliance by increasing enforcement and imposing greater
      fines for breach. Beyond these facts it's simply the right thing to do. At
      Thunderdome we strongly believe that your data privacy is very important
      and we already have solid security and privacy practices in place that go
      beyond the requirements of this new regulation.
    </p>

    <h2 class="py-2 text-xl">
      Individual Data Subject's Rights - Data Access, Portability and Deletion
    </h2>
    <p class="py-4">
      We are committed to helping our customers meet the data subject rights
      requirements of GDPR. Thunderdome processes or stores all personal data in
      fully vetted, DPA compliant vendors. We do store all conversation and
      personal data for up to 6 years unless your account is deleted. In which
      case, we dispose of all data in accordance with our Terms of Service and
      Privacy Policy, but we will not hold it longer than 60 days.
    </p>
    <p class="py-4">
      We are aware that if you are working with EU customers, you need to be
      able to provide them with the ability to access, update, retrieve and
      remove personal data. We got you! We've been set up as self service from
      the start and have always given you access to your data and your customers
      data. Our customer support team is here for you to answer any questions
      you might have about working with the API.
    </p>
    <h2 class="py-2 text-xl">California Residents</h2>
    <p class="py-4">
      The California Consumer Privacy Act (CCPA) requires us to disclose
      categories of Personal Information we collect and how we use it, the
      categories of sources from whom we collect Personal Information, and the
      third parties with whom we share it, which we have explained above.
    </p>
    <p class="py-4">
      We are also required to communicate information about rights California
      residents have under California law. You may exercise the following
      rights:
    </p>
    <ul class="space-y-1 list-disc list-inside">
      <li>
        Right to Know and Access. You may submit a verifiable request for
        information regarding the: (1) categories of Personal Information we
        collect, use, or share; (2) purposes for which categories of Personal
        Information are collected or used by us; (3) categories of sources from
        which we collect Personal Information; and (4) specific pieces of
        Personal Information we have collected about you.
      </li>
      <li>
        Right to Equal Service. We will not discriminate against you if you
        exercise your privacy rights.
      </li>
      <li>
        Right to Delete. You may submit a verifiable request to close your
        account and we will delete Personal Information about you that we have
        collected.
      </li>
      <li>
        Request that a business that sells a consumer's personal data, not sell
        the consumer's personal data.
      </li>
    </ul>
    <p class="py-4">
      If you make a request, we have one month to respond to you. If you would
      like to exercise any of these rights, please contact us.
    </p>
    <p class="py-4">We do not sell the Personal Information of our users.</p>
    <p class="py-4">
      For more information about these rights, please contact us.
    </p>
    <h2 class="py-2 text-xl">
      California Online Privacy Protection Act (CalOPPA)
    </h2>
    <p class="py-4">
      CalOPPA requires us to disclose categories of Personal Information we
      collect and how we use it, the categories of sources from whom we collect
      Personal Information, and the third parties with whom we share it, which
      we have explained above.
    </p>
    <p class="py-4">CalOPPA users have the following rights:</p>
    <ul class="space-y-1 list-disc list-inside">
      <li>
        Right to Know and Access. You may submit a verifiable request for
        information regarding the: (1) categories of Personal Information we
        collect, use, or share; (2) purposes for which categories of Personal
        Information are collected or used by us; (3) categories of sources from
        which we collect Personal Information; and (4) specific pieces of
        Personal Information we have collected about you.
      </li>
      <li>
        Right to Equal Service. We will not discriminate against you if you
        exercise your privacy rights.
      </li>
      <li>
        Right to Delete. You may submit a verifiable request to close your
        account and we will delete Personal Information about you that we have
        collected.
      </li>
      <li>
        Right to request that a business that sells a consumer's personal data,
        not sell the consumer's personal data.
      </li>
    </ul>
    <p class="py-4">
      If you make a request, we have one month to respond to you. If you would
      like to exercise any of these rights, please contact us.
    </p>
    <p class="py-4">We do not sell the Personal Information of our users.</p>
    <p class="py-4">
      For more information about these rights, please contact us.
    </p>
    <h2 class="py-2 text-xl">Contact Us</h2>
    <p class="py-4">Don't hesitate to contact us if you have any questions.</p>
    <ul class="space-y-1 list-disc list-inside">
      <li>
        Via this Link: <a
          href="https://thunderdome.dev/support"
          class="text-cyan-600 dark:text-cyan-400 hover:text-thunder-yellow"
          >https://thunderdome.dev/support</a
        >
      </li>
    </ul>
  </div>
</PageLayout>
